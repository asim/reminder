#!/usr/bin/env python3
"""
Scrape hadith data from sunnah.com with Arabic and English text.
"""

import json
import re
import time
import sys
import requests
from bs4 import BeautifulSoup
from pathlib import Path

# Collections to scrape
COLLECTIONS = {
    'bukhari': {'name': 'Sahih al-Bukhari', 'arabic': 'صحيح البخاري'},
    'muslim': {'name': 'Sahih Muslim', 'arabic': 'صحيح مسلم'},
    'nawawi40': {'name': "An-Nawawi's 40 Hadith", 'arabic': 'الأربعون النووية'},
}

def clean_book_name(name):
    """Clean up book name by separating number and removing duplicate Arabic."""
    # Remove leading numbers
    name = re.sub(r'^\d+', '', name)
    # Split on Arabic characters if English+Arabic mixed
    parts = re.split(r'([\u0600-\u06FF]+)', name)
    english = ''.join(p for p in parts if not re.search(r'[\u0600-\u06FF]', p)).strip()
    return english if english else name

def get_collection_books(collection):
    """Get list of books for a collection."""
    url = f"https://sunnah.com/{collection}"
    resp = requests.get(url, timeout=30)
    soup = BeautifulSoup(resp.text, 'html.parser')
    
    books = []
    # Find book links - they're typically in a list
    for link in soup.select(f'a[href^="/{collection}/"]'):
        href = link.get('href', '')
        match = re.match(rf'/{collection}/(\d+)$', href)
        if match:
            book_num = int(match.group(1))
            name = link.get_text(strip=True)
            books.append({'number': book_num, 'name': clean_book_name(name)})
    
    # Deduplicate and sort
    seen = set()
    unique_books = []
    for b in sorted(books, key=lambda x: x['number']):
        if b['number'] not in seen:
            seen.add(b['number'])
            unique_books.append(b)
    
    return unique_books

def get_book_hadiths(collection, book_num):
    """Get all hadiths from a book."""
    url = f"https://sunnah.com/{collection}/{book_num}"
    resp = requests.get(url, timeout=30)
    soup = BeautifulSoup(resp.text, 'html.parser')
    
    hadiths = []
    seen_numbers = set()
    
    # Find all hadith containers
    for container in soup.select('.hadithTextContainers, .actualHadithContainer'):
        hadith = {}
        
        # Get hadith number from reference
        ref = container.select_one('.hadith_reference_sticky, .hadith_reference')
        if ref:
            ref_text = ref.get_text(strip=True)
            num_match = re.search(r'(\d+)$', ref_text)
            if num_match:
                hadith['number'] = int(num_match.group(1))
        
        # Skip duplicates
        if hadith.get('number') in seen_numbers:
            continue
        if hadith.get('number'):
            seen_numbers.add(hadith['number'])
        
        # Get narrator
        narrated = container.select_one('.hadith_narrated')
        if narrated:
            hadith['narrator'] = narrated.get_text(strip=True)
        
        # Get English text
        english = container.select_one('.text_details')
        if english:
            # Clean up whitespace
            text = ' '.join(english.get_text().split())
            hadith['english'] = text
        
        # Get Arabic text (the actual hadith, not the chain)
        arabic_text = container.select_one('.arabic_text_details')
        if arabic_text:
            hadith['arabic'] = arabic_text.get_text(strip=True)
        
        # Get Arabic chain (sanad) - optional
        arabic_sanad = container.select_one('.arabic_sanad')
        if arabic_sanad:
            hadith['chain'] = arabic_sanad.get_text(strip=True)
        
        # Only add if we have content
        if hadith.get('english') or hadith.get('arabic'):
            hadiths.append(hadith)
    
    return hadiths

def scrape_collection(collection, info, output_dir):
    """Scrape an entire collection."""
    print(f"Scraping {info['name']}...")
    
    books = get_collection_books(collection)
    print(f"  Found {len(books)} books")
    
    result = {
        'name': info['name'],
        'arabic': info['arabic'],
        'collection': collection,
        'books': []
    }
    
    total_hadiths = 0
    
    for book in books:
        print(f"  Book {book['number']}: {book['name']}", end='', flush=True)
        time.sleep(0.5)  # Be nice to the server
        
        try:
            hadiths = get_book_hadiths(collection, book['number'])
            print(f" - {len(hadiths)} hadiths")
            total_hadiths += len(hadiths)
            
            result['books'].append({
                'number': book['number'],
                'name': book['name'],
                'hadiths': hadiths
            })
        except Exception as e:
            print(f" - ERROR: {e}")
    
    print(f"  Total: {total_hadiths} hadiths")
    
    # Save to file
    output_file = output_dir / f"{collection}.json"
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(result, f, ensure_ascii=False, indent=2)
    
    print(f"  Saved to {output_file}")
    return result

def main():
    output_dir = Path('hadith/data_new')
    output_dir.mkdir(exist_ok=True)
    
    # Scrape specified collections or all
    collections_to_scrape = sys.argv[1:] if len(sys.argv) > 1 else list(COLLECTIONS.keys())
    
    for collection in collections_to_scrape:
        if collection not in COLLECTIONS:
            print(f"Unknown collection: {collection}")
            continue
        scrape_collection(collection, COLLECTIONS[collection], output_dir)
        print()

if __name__ == '__main__':
    main()
