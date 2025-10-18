import requests
import json
import base64
import os
from pathlib import Path

# Configuration
BASE_URL = "http://localhost:8080"  # Change this to your server URL
CREATE_POST_ENDPOINT = f"{BASE_URL}/api/v1/posts"

# Test data
test_posts = [
    {
        "name": "Create Post - Electronics",
        "body": {
            "creatorId": "firebase-user-001",
            "creatorMail": "seller1@example.com",
            "title": "MacBook Pro 2020 for Sale",
            "description": "Excellent condition MacBook Pro with M1 chip",
            "tags": ["electronics", "laptop", "apple"],
            "text": "Selling my MacBook Pro 2020 with M1 chip, 16GB RAM, 512GB SSD. Barely used, comes with original charger and box. Perfect for students or professionals.",
            "payPalMail": "seller1@example.com",
            "images": ["./test-images/macbook.jpg"]
        }
    },
    {
        "name": "Create Post - Furniture",
        "body": {
            "creatorId": "firebase-user-002",
            "creatorMail": "furniture.lover@example.com",
            "title": "Vintage Wooden Desk",
            "description": "Beautiful handcrafted oak desk from the 1960s",
            "tags": ["furniture", "vintage", "desk", "wood"],
            "text": "Stunning vintage oak desk in excellent condition. Features 3 drawers, brass handles, and a leather writing surface. Dimensions: 120cm x 60cm x 75cm.",
            "payPalMail": "furniture.lover@example.com",
            "images": ["./test-images/desk1.jpg", "./test-images/desk2.jpg"]
        }
    },
    {
        "name": "Create Post - Camera",
        "body": {
            "creatorId": "firebase-user-003",
            "creatorMail": "photo.pro@example.com",
            "title": "Canon EOS R5 Camera Body",
            "description": "Professional mirrorless camera, like new",
            "tags": ["electronics", "camera", "photography", "canon"],
            "text": "Canon EOS R5 body only, purchased 6 months ago. Shutter count under 5000. Includes original packaging, battery, charger, and strap. Perfect for professional photographers.",
            "payPalMail": "photo.pro@example.com",
            "images": ["./test-images/camera.jpg"]
        }
    },
    {
        "name": "Create Post - Bicycle",
        "body": {
            "creatorId": "firebase-user-004",
            "creatorMail": "bike.enthusiast@example.com",
            "title": "Mountain Bike - Trek X-Caliber 9",
            "description": "2022 model, carbon frame, excellent condition",
            "tags": ["sports", "bicycle", "outdoor", "trek"],
            "text": "Trek X-Caliber 9 mountain bike from 2022. Carbon frame, 29-inch wheels, Shimano Deore XT components. Recently serviced, new tires. Size: Large (suitable for 180-195cm height).",
            "payPalMail": "bike.enthusiast@example.com",
            "images": ["./test-images/bike1.jpg", "./test-images/bike2.jpg", "./test-images/bike3.jpg"]
        }
    },
    {
        "name": "Create Post - Books",
        "body": {
            "creatorId": "firebase-user-005",
            "creatorMail": "bookworm@example.com",
            "title": "Computer Science Textbook Collection",
            "description": "Set of 12 CS textbooks for university students",
            "tags": ["books", "education", "computer-science"],
            "text": "Selling my computer science textbook collection from university. Includes algorithms, data structures, operating systems, and more. All books in good condition with minimal highlighting. Perfect for CS students.",
            "payPalMail": "bookworm@example.com",
            "images": ["./test-images/books.jpg"]
        }
    },
    {
        "name": "Create Post - Guitar",
        "body": {
            "creatorId": "firebase-user-006",
            "creatorMail": "musician123@example.com",
            "title": "Fender Stratocaster Electric Guitar",
            "description": "American Professional Series, Sunburst finish",
            "tags": ["music", "guitar", "fender", "instruments"],
            "text": "Fender American Professional Stratocaster in 3-Color Sunburst. Rosewood fingerboard, V-Mod pickups. Played regularly but well maintained. Comes with hardshell case and strap.",
            "payPalMail": "musician123@example.com",
            "images": ["./test-images/guitar1.jpg", "./test-images/guitar2.jpg"]
        }
    },
    {
        "name": "Create Post - Gaming Console",
        "body": {
            "creatorId": "firebase-user-007",
            "creatorMail": "gamer.paradise@example.com",
            "title": "PlayStation 5 with 2 Controllers",
            "description": "PS5 Disc Edition + extra controller + 5 games",
            "tags": ["gaming", "console", "playstation", "electronics"],
            "text": "PlayStation 5 Disc Edition in perfect condition. Includes 2 DualSense controllers, 5 popular games (God of War, Spider-Man, etc.), and all cables. Barely used, smoke-free home.",
            "payPalMail": "gamer.paradise@example.com",
            "images": ["./test-images/ps5.jpg"]
        }
    },
    {
        "name": "Create Post - Coffee Machine",
        "body": {
            "creatorId": "firebase-user-008",
            "creatorMail": "coffee.addict@example.com",
            "title": "De'Longhi Espresso Machine",
            "description": "Professional-grade espresso and cappuccino maker",
            "tags": ["appliances", "coffee", "kitchen", "delonghi"],
            "text": "De'Longhi La Specialista Prestigio espresso machine. Dual heating system, built-in grinder, automatic milk frother. Used for 1 year, regularly cleaned and maintained. Perfect for coffee lovers.",
            "payPalMail": "coffee.addict@example.com",
            "images": ["./test-images/espresso1.jpg", "./test-images/espresso2.jpg"]
        }
    },
    {
        "name": "Create Post - Running Shoes",
        "body": {
            "creatorId": "firebase-user-009",
            "creatorMail": "runner.pro@example.com",
            "title": "Nike Air Zoom Pegasus 40 - Size 10.5",
            "description": "Brand new, never worn, with box and tags",
            "tags": ["shoes", "sports", "running", "nike"],
            "text": "Nike Air Zoom Pegasus 40 running shoes in size US 10.5 / EU 44.5. Received as gift but wrong size. Never worn, still in original box with all tags. Color: Black/White.",
            "payPalMail": "runner.pro@example.com",
            "images": ["./test-images/shoes.jpg"]
        }
    },
    {
        "name": "Create Post - Smart Watch",
        "body": {
            "creatorId": "firebase-user-010",
            "creatorMail": "tech.savvy@example.com",
            "title": "Apple Watch Series 8 - 45mm GPS",
            "description": "Midnight aluminum with Sport Band, like new",
            "tags": ["electronics", "smartwatch", "apple", "fitness"],
            "text": "Apple Watch Series 8 in midnight aluminum, 45mm case with GPS. Includes original midnight Sport Band plus an extra braided loop. Screen protector applied since day one. Battery health at 100%. Includes charging cable and box.",
            "payPalMail": "tech.savvy@example.com",
            "images": ["./test-images/watch1.jpg", "./test-images/watch2.jpg", "./test-images/watch3.jpg"]
        }
    }
]

def load_image_as_base64(image_path):
    """Load an image file and convert it to base64 (without prefix)"""
    try:
        with open(image_path, 'rb') as image_file:
            encoded_string = base64.b64encode(image_file.read()).decode('utf-8')
            return encoded_string
    except FileNotFoundError:
        print(f"Warning: Image file '{image_path}' not found. Using placeholder.")
        # Return a tiny 1x1 transparent PNG as placeholder
        return "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg=="

def process_images(images_list):
    """Process image paths and convert to base64"""
    base64_images = []
    for image_path in images_list:
        base64_image = load_image_as_base64(image_path)
        base64_images.append(base64_image)
    return base64_images

def create_post(post_data):
    """Send POST request to create a post"""
    # Process images
    body = post_data["body"].copy()
    body["images"] = process_images(body["images"])
    
    headers = {
        "Content-Type": "application/json"
    }
    
    try:
        response = requests.post(
            CREATE_POST_ENDPOINT,
            json=body,
            headers=headers,
            timeout=30
        )
        
        return {
            "name": post_data["name"],
            "status_code": response.status_code,
            "response": response.text,
            "success": response.status_code == 200
        }
    except requests.exceptions.RequestException as e:
        return {
            "name": post_data["name"],
            "status_code": None,
            "response": str(e),
            "success": False
        }

def main():
    print(f"Starting tests against {BASE_URL}")
    print("=" * 80)
    
    results = []
    success_count = 0
    
    for i, post_data in enumerate(test_posts, 1):
        print(f"\nTest {i}/{len(test_posts)}: {post_data['name']}")
        print("-" * 80)
        
        result = create_post(post_data)
        results.append(result)
        
        if result["success"]:
            success_count += 1
            print(f"✓ SUCCESS - Status: {result['status_code']}")
            print(f"Response: {result['response']}")
        else:
            print(f"✗ FAILED - Status: {result['status_code']}")
            print(f"Response: {result['response']}")
    
    # Summary
    print("\n" + "=" * 80)
    print(f"SUMMARY: {success_count}/{len(test_posts)} tests passed")
    print("=" * 80)
    
    # Detailed results
    print("\nDetailed Results:")
    for i, result in enumerate(results, 1):
        status = "✓" if result["success"] else "✗"
        print(f"{status} {i}. {result['name']} - Status: {result['status_code']}")

if __name__ == "__main__":
    main()