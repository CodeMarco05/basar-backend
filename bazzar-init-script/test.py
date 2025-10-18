import requests
import json
import base64
import os
from pathlib import Path

# Configuration
BASE_URL = "http://localhost:42001"  # Change this to your server URL
CREATE_POST_ENDPOINT = f"{BASE_URL}/api/v1/posts"

# Test data
test_posts = [
    {
    "name": "Treasure Hoard",
    "body": {
        "creatorId": "firebase-user-011",
        "creatorMail": "totally.not.a.dragon@example.com",
        "title": "Genuine Dragon's Treasure - Must Sell ASAP",
        "description": "Relocating to smaller cave, need to downsize collection",
        "tags": ["gold", "treasure", "collectibles", "definitely-real"],
        "text": "After 500 years of hoarding, I've finally admitted I have a problem. Selling my gold pile (approximately 3 tons). NO LOWBALLERS - I know what I have! Cash only. Serious inquiries only - if you're a knight, don't bother. Pick-up only from mountain cave. May throw in cursed amulet for free if you take it all. P.S. - Ignore the scorch marks on the coins.",
        "payPalMail": "totally.not.a.dragon@example.com",
        "images": ["./test-images/gold-pile.jpg"]
    }
    },
    {
    "name": "Uni zu verkaufen",
    "body": {
        "creatorId": "firebase-user-012",
        "creatorMail": "dr.sauer@nordakademie.de",
        "title": "Komplette Universität zu verkaufen - Gebraucht aber funktional",
        "description": "NAK Elmshorn - Wegen Budgetkürzungen muss leider alles weg",
        "tags": ["immobilien", "bildung", "schnäppchen", "selbstabholer"],
        "text": "Verkaufe hier meine gebrauchte Universität in Elmshorn. Zustand: Stark frequentiert, aber noch voll funktionsfähig. Enthält mehrere Hörsäle (Sitze teilweise durchgesessen), Mensa (Kaffeemaschine inkl.), Bibliothek (Bücher gegen Aufpreis), und diverse Büros (IKEA-Möbel bleiben drin). Perfekt für den ambitionierten Bildungsunternehmer oder als Studentenwohnheim-Conversion-Projekt. ACHTUNG: Studenten sind NICHT im Lieferumfang enthalten und müssen separat erworben werden. Abholung nur komplett, kein Teilverkauf einzelner Räume. VB 50 Mio. € - Verhandlungsbasis bei Barzahlung. PS: Parkplatzproblem wird kostenlos mitgeliefert.",
        "payPalMail": "prof.dr.sparsam@nordakademie.de",
        "images": ["./test-images/nak-elmshorn.jpg"]
    }
    },
{
    "name": "Seminar Platz abzugeben",
    "body": {
        "creatorId": "firebase-user-013",
        "creatorMail": "verzweifelter.student@nordakademie.de",
        "title": "DRINGEND: Seminar-Platz abzugeben - Habe keine Zeit mehr",
        "description": "Seminar 'Die Big Five der beruflichen Umgangsformen' - Kann leider nicht teilnehmen",
        "tags": ["bildung", "seminar", "student", "nordakademie", "verzweifelt"],
        "text": "Verkaufe meinen hart erkämpften Seminar-Platz für 'Die Big Five der beruflichen Umgangsformen' am kommenden Wochenende.",
        "payPalMail": "verzweifelter.student@nordakademie.de",
        "images": ["./test-images/seminar-banner.png"]
    }
},
    {
        "name": "Business Seminar Platz",
        "body": {
            "creatorId": "firebase-user-014",
            "creatorMail": "startup.dreams@nordakademie.de",
            "title": "Seminar-Platz: Von der Businessidee zum Businessplan",
            "description": "Seminar am Wochenende - Muss leider absagen",
            "tags": ["seminar", "business", "nordakademie", "startup", "bildung"],
            "text": "Biete Seminar-Platz für 'Von der Businessidee zum Businessplan' am kommenden Wochenende. Hatte mir vorgenommen, endlich meine geniale App-Idee umzusetzen, aber Netflix hat eine neue Serie rausgebracht. Vielleicht nächstes Semester. Der Seminarleiter soll super sein! Perfekt für angehende Gründer oder alle, die endlich ihre Ideen zu Papier bringen wollen. VB 20€ oder ein gutes Pitch-Deck als Tausch.",
            "payPalMail": "startup.dreams@nordakademie.de",
            "images": ["./test-images/tony-seminar.jpg"]
        }
    },
    {
        "name": "Transferleistung ghostwriter",
        "body": {
            "creatorId": "firebase-user-014",
            "creatorMail": "desperate.student@nordakademie.de",
            "title": "Transferleistung Ghostwriter gesucht - DRINGEND!!!",
            "description": "Suche jemanden, der meine TL schreibt - Deadline in 2 Wochen",
            "tags": ["transferleistung", "ghostwriter", "nordakademie", "studium", "verzweifelt"],
            "text": "Biete großzügige Bezahlung für jemanden, der meine Transferleistung zum Thema 'Digitalisierung in mittelständischen Unternehmen' schreibt. Hatte ursprünglich vor, das selbst zu machen, aber habe die letzten 3 Monate nur an meiner Work-Life-Balance gearbeitet (hauptsächlich Life, weniger Work). Anforderungen: 15-20 Seiten, wissenschaftlicher Stil, Minimum 10 Quellen. Mein Praxisunternehmen weiß von nichts und soll auch nichts erfahren! Biete: 300€ VB, ewige Dankbarkeit, oder meine Seminarplätze für die nächsten 2 Semester. WICHTIG: Plagiatsprüfung muss bestanden werden! Bei Interesse bitte diskrete Kontaktaufnahme. P.S.: Noten zwischen 2,0-3,0 sind völlig ok, will nicht auffallen. 😅",
            "payPalMail": "desperate.student@nordakademie.de",
            "images": ["./test-images/transferleistung-loewe.jpeg"]
        }
    },
    {
        "name": "Emotionale Intelligenz Seminar Platz",
        "body": {
            "creatorId": "firebase-user-014",
            "creatorMail": "startup.dreams@nordakademie.de",
            "title": "Emotionale Intelligenz Seminar Platz",
            "description": "Seminar am Wochenende - Muss leider absagen",
            "tags": ["seminar", "wochenende", "nordakademie"],
            "text": "Biete Seminar-Platz für 'Emotionale Intelligenz' am kommenden Wochenende.",
            "payPalMail": "startup.dreams@nordakademie.de",
            "images": ["./test-images/wochenende-an-nak.jpg"]
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