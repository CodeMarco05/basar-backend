"""
Scraper for dasistschoen.de using Selenium to handle JavaScript rendering
"""

from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.common.exceptions import TimeoutException, NoSuchElementException
import time
import json
from datetime import datetime
import base64
import requests
from urllib.parse import urljoin
import os
from pymongo import MongoClient, errors
from dotenv import load_dotenv
from PIL import Image
import io

# Load environment variables from .env file
load_dotenv()


class DasistschoenSeleniumScraper:
    def __init__(self, headless=True, use_mongodb=True):
        self.base_url = 'https://dasistschoen.de'
        self.all_images = []
        self.headless = headless
        self.driver = None
        self.use_mongodb = use_mongodb
        self.mongo_client = None
        self.mongo_db = None
        self.mongo_collection = None
        
        if self.use_mongodb:
            self.setup_mongodb()
        
    def setup_mongodb(self):
        """Initialize MongoDB connection"""
        try:
            mongodb_uri = os.getenv('SCHABERNAK_MONGODB_URI', 'mongodb://admin:dasistschoen123@localhost:27017/')
            mongodb_database = os.getenv('SCHABERNAK_MONGODB_DATABASE', 'dasistschoen')
            mongodb_collection = os.getenv('SCHABERNAK_MONGODB_COLLECTION', 'memes')
            
            print(f"Connecting to MongoDB at {mongodb_uri}...")
            self.mongo_client = MongoClient(mongodb_uri, serverSelectionTimeoutMS=5000)
            
            # Test connection
            self.mongo_client.server_info()
            
            self.mongo_db = self.mongo_client[mongodb_database]
            self.mongo_collection = self.mongo_db[mongodb_collection]
            
            # Create index on id field for faster lookups and prevent duplicates
            self.mongo_collection.create_index("id", unique=True)
            
            print(f"✓ Connected to MongoDB database '{mongodb_database}', collection '{mongodb_collection}'")
            
        except errors.ServerSelectionTimeoutError:
            print("✗ Could not connect to MongoDB. Falling back to JSON-only mode.")
            self.use_mongodb = False
        except Exception as e:
            print(f"✗ MongoDB setup error: {e}. Falling back to JSON-only mode.")
            self.use_mongodb = False
    
    def setup_driver(self):
        """Initialize the Chrome WebDriver"""
        options = webdriver.ChromeOptions()
        if self.headless:
            options.add_argument('--headless')
        options.add_argument('--no-sandbox')
        options.add_argument('--disable-dev-shm-usage')
        options.add_argument('user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36')
        
        self.driver = webdriver.Chrome(options=options)
        self.driver.implicitly_wait(10)
        
    def scrape_all(self, max_clicks=None):
        """
        Scrape all images by clicking the 'Mehr laden' button
        
        Args:
            max_clicks: Maximum number of times to click "Load More" button (None = unlimited)
        """
        self.setup_driver()
        
        try:
            print("Loading homepage...")
            self.driver.get(self.base_url)
            
            # Wait for initial content to load
            WebDriverWait(self.driver, 10).until(
                EC.presence_of_element_located((By.CLASS_NAME, "image-container"))
            )
            print("Page loaded successfully!")
            
            # First, click "Load More" button the specified number of times
            clicks = 0
            while True:
                if max_clicks and clicks >= max_clicks:
                    print(f"Reached maximum clicks limit: {max_clicks}")
                    break
                
                # Try to find and click the "Mehr laden" button
                try:
                    load_more_button = self.driver.find_element(By.ID, "btn-more")
                    
                    # Check if button is still visible
                    if not load_more_button.is_displayed():
                        print("Load more button is no longer visible. All images loaded!")
                        break
                    
                    # Scroll to button
                    self.driver.execute_script("arguments[0].scrollIntoView(true);", load_more_button)
                    time.sleep(0.5)
                    
                    # Click the button
                    print(f"Clicking 'Mehr laden' button (click #{clicks + 1})...")
                    load_more_button.click()
                    clicks += 1
                    
                    # Wait for new content to load
                    time.sleep(2)
                    
                except NoSuchElementException:
                    print("No more 'Mehr laden' button found. All images loaded!")
                    break
                except Exception as e:
                    print(f"Error clicking button: {e}")
                    break
            
            # After all clicks, extract all images at once
            print(f"\nFinished clicking. Now extracting all image data...")
            self.all_images = self.extract_images()
            print(f"✓ Total images extracted: {len(self.all_images)}")
            
        finally:
            self.driver.quit()
        
        return self.all_images
    
    def download_image_as_base64(self, image_url, quality=85):
        """Download image, compress it, and convert to base64

        Args:
            image_url: URL of the image to download
            quality: JPEG quality (1-100, higher is better quality)
        """
        try:
            # Make absolute URL if relative
            full_url = urljoin(self.base_url, image_url)

            # Download image
            response = requests.get(full_url, timeout=10)
            response.raise_for_status()

            # Open image with Pillow
            img = Image.open(io.BytesIO(response.content))

            # Convert RGBA to RGB if necessary (for JPEG compatibility)
            if img.mode in ('RGBA', 'LA', 'P'):
                background = Image.new('RGB', img.size, (255, 255, 255))
                if img.mode == 'P':
                    img = img.convert('RGBA')
                background.paste(img, mask=img.split()[-1] if img.mode in ('RGBA', 'LA') else None)
                img = background

            # Compress and save to bytes
            output_buffer = io.BytesIO()
            img.save(output_buffer, format='JPEG', quality=quality, optimize=True)
            compressed_data = output_buffer.getvalue()

            # Convert to base64
            base64_data = base64.b64encode(compressed_data).decode('utf-8')

            # Return as data URI (always JPEG after compression)
            return f"data:image/jpeg;base64,{base64_data}"

        except Exception as e:
            print(f"Error downloading/compressing image {image_url}: {e}")
            return None
    
    def extract_images(self):
        """Extract image data from currently loaded page"""
        images = []
        
        # Find all image containers
        image_containers = self.driver.find_elements(By.CLASS_NAME, "image-container")
        
        for idx, container in enumerate(image_containers, 1):
            try:
                image_data = {}
                
                # Extract image URL
                img_element = container.find_element(By.CSS_SELECTOR, ".the-image img")
                image_url = img_element.get_attribute('src')
                image_data['image_url'] = image_url
                image_data['alt_text'] = img_element.get_attribute('alt')
                
                # Download and convert image to base64
                print(f"Downloading image {idx}/{len(image_containers)}...", end=' ')
                base64_image = self.download_image_as_base64(image_url)
                if base64_image:
                    image_data['image_base64'] = base64_image
                    print("✓")
                else:
                    print("✗")
                
                # Extract title
                title_element = container.find_element(By.CSS_SELECTOR, ".title .title-text")
                image_data['title'] = title_element.text.strip()
                
                # Extract permalink
                try:
                    permalink = container.find_element(By.CLASS_NAME, "image-permalink")
                    image_data['permalink'] = permalink.get_attribute('href')
                    # Extract ID from permalink
                    if '/image/' in image_data['permalink']:
                        image_data['id'] = image_data['permalink'].split('/image/')[-1]
                except:
                    pass
                
                # Extract likes
                try:
                    likes_element = container.find_element(By.CLASS_NAME, "likebutton")
                    likes_text = likes_element.text.strip()
                    # Extract number from text like "♥ 5"
                    image_data['likes'] = ''.join(filter(str.isdigit, likes_text))
                except:
                    image_data['likes'] = '0'
                
                # Extract timestamp
                try:
                    timestamp_element = container.find_element(By.CLASS_NAME, "timestamp")
                    image_data['timestamp'] = timestamp_element.text.strip()
                except:
                    pass
                
                # Extract user
                try:
                    user_element = container.find_element(By.CLASS_NAME, "user-name")
                    image_data['user'] = user_element.text.strip()
                except:
                    pass
                
                # Extract tags
                try:
                    tags_elements = container.find_elements(By.CSS_SELECTOR, ".tags a")
                    image_data['tags'] = [tag.text.strip() for tag in tags_elements]
                except:
                    image_data['tags'] = []
                
                # Extract comments
                try:
                    comments_section = container.find_element(By.CLASS_NAME, "comments")
                    comment_elements = comments_section.find_elements(By.CLASS_NAME, "image-comment")
                    image_data['comments'] = []
                    for comment in comment_elements:
                        if 'no-comments' not in comment.get_attribute('class'):
                            image_data['comments'].append(comment.text.strip())
                except:
                    image_data['comments'] = []
                
                images.append(image_data)
                
            except Exception as e:
                print(f"Error extracting image data: {e}")
                continue
        
        return images
    
    def save_results(self, filename=None):
        """Save scraped images to JSON file and MongoDB"""
        if not filename:
            timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
            # Save to /app/output in Docker, or current directory otherwise
            output_dir = '/app/output' if os.path.exists('/app/output') else './output'
            os.makedirs(output_dir, exist_ok=True)
            filename = f'{output_dir}/dasistschoen_memes_{timestamp}.json'
        
        # Save to JSON file
        with open(filename, 'w', encoding='utf-8') as f:
            json.dump(self.all_images, f, indent=2, ensure_ascii=False)
        
        print(f"\n✓ Saved {len(self.all_images)} images to {filename}")
        
        # Save to MongoDB
        if self.use_mongodb and self.mongo_collection is not None:
            self.save_to_mongodb()
        
        return filename
    
    def save_to_mongodb(self):
        """Save or update images in MongoDB"""
        if not self.all_images:
            print("No images to save to MongoDB")
            return
        
        print(f"\nSaving to MongoDB...")
        inserted = 0
        updated = 0
        unchanged = 0
        errors_count = 0
        
        for image_data in self.all_images:
            try:
                # Add scraped timestamp
                image_data['scraped_at'] = datetime.now()
                
                # Use upsert to insert or update based on image ID
                if 'id' in image_data:
                    result = self.mongo_collection.update_one(
                        {'id': image_data['id']},
                        {'$set': image_data},
                        upsert=True
                    )
                    
                    if result.upserted_id:
                        inserted += 1
                    elif result.modified_count > 0:
                        updated += 1
                    else:
                        unchanged += 1
                else:
                    # If no ID, just insert (may create duplicates)
                    self.mongo_collection.insert_one(image_data)
                    inserted += 1
                    
            except Exception as e:
                print(f"Error saving image to MongoDB: {e}")
                errors_count += 1
        
        print(f"✓ MongoDB: {inserted} new, {updated} updated, {unchanged} unchanged, {errors_count} errors")
        print(f"✓ Total documents in collection: {self.mongo_collection.count_documents({})}")
    
    def print_summary(self):
        """Print scraping summary"""
        print("\n" + "="*50)
        print("SCRAPING SUMMARY")
        print("="*50)
        print(f"Total images scraped: {len(self.all_images)}")
        
        if self.all_images:
            # Show sample of first image (convert datetime for JSON serialization)
            print("\nSample image data:")
            sample = self.all_images[0].copy()
            if 'scraped_at' in sample:
                sample['scraped_at'] = sample['scraped_at'].isoformat() if hasattr(sample['scraped_at'], 'isoformat') else str(sample['scraped_at'])
            print(json.dumps(sample, indent=2, ensure_ascii=False))


def main():
    # Use headless mode in Docker, allow non-headless locally
    is_docker = os.path.exists('/app/output')
    
    # Check if MongoDB should be used (default: True)
    use_mongodb = os.getenv('USE_MONGODB', 'true').lower() != 'false'
    
    # Get scrape interval (in seconds) - if set, run continuously
    scrape_interval = os.getenv('MEMES_SCRAPE_INTERVAL')
    if scrape_interval:
        try:
            scrape_interval = int(scrape_interval)
            print(f"Continuous scraping mode enabled: running every {scrape_interval} seconds")
        except ValueError:
            print(f"Invalid MEMES_SCRAPE_INTERVAL value: {scrape_interval}. Using single-run mode.")
            scrape_interval = None
    
    # Get max_clicks setting
    max_clicks_env = os.getenv('MAX_CLICKS')
    max_clicks = int(max_clicks_env) if max_clicks_env else None
    
    # Continuous scraping mode
    if scrape_interval:
        print("Starting Selenium scraper for dasistschoen.de in CONTINUOUS MODE")
        print("="*50)
        
        run_count = 0
        while True:
            run_count += 1
            print(f"\n{'='*50}")
            print(f"SCRAPE RUN #{run_count} - {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
            print("="*50)
            
            try:
                # Create a new scraper instance for each run
                scraper = DasistschoenSeleniumScraper(headless=True, use_mongodb=use_mongodb)
                
                images = scraper.scrape_all(max_clicks=max_clicks)
                
                if images:
                    scraper.save_results()
                    scraper.print_summary()
                else:
                    print("\nNo images were scraped.")
                
                # Close MongoDB connection
                if scraper.mongo_client:
                    scraper.mongo_client.close()
                    
            except Exception as e:
                print(f"\n✗ Error during scrape run #{run_count}: {e}")
                import traceback
                traceback.print_exc()
            
            print(f"\n✓ Scrape run #{run_count} completed. Waiting {scrape_interval} seconds until next run...")
            print(f"Next run at: {(datetime.now() + __import__('datetime').timedelta(seconds=scrape_interval)).strftime('%Y-%m-%d %H:%M:%S')}")
            
            time.sleep(scrape_interval)
    
    # Single run mode
    else:
        print("Starting Selenium scraper for dasistschoen.de in SINGLE-RUN MODE")
        print("="*50)
        
        scraper = DasistschoenSeleniumScraper(headless=True, use_mongodb=use_mongodb)
        
        images = scraper.scrape_all(max_clicks=max_clicks)
        
        if images:
            scraper.save_results()
            scraper.print_summary()
        else:
            print("\nNo images were scraped.")
        
        # Close MongoDB connection
        if scraper.mongo_client:
            scraper.mongo_client.close()
            print("\n✓ MongoDB connection closed")


if __name__ == "__main__":
    main()
