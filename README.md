# 🕷️ Simple Web Crawler for Adidas Japan

**Target URL:** [https://shop.adidas.jp/men/](https://shop.adidas.jp/men/)

This project is my first attempt at building a web crawler in Go. It's a bit messy in places, but it was an incredibly valuable learning experience.

---

## 📚 What I Learned

- Handling bot detection and anti-scraping mechanisms.
- Extracting data from both HTML and embedded JavaScript.
- Working with multiple Go-based scraping tools:
    - Raw `net/http` requests
    - [`colly`](https://github.com/gocolly/colly) (Go crawling framework)
    - [`rod`](https://github.com/go-rod/rod) for headless browser automation
    - [`chromedp`](https://github.com/chromedp/chromedp) (which finally worked — sort of!)

---

## 🚧 Challenges Faced

- **Bot Detection**:  
  Most requests triggered Adidas’s bot protection, returning obfuscated JavaScript instead of HTML.

- **Headless Browsers**:  
  Tools like Rod and Colly failed to bypass detection. Even Chromedp only worked intermittently with carefully timed delays.

- **API Discovery**:  
  I managed to discover some of Adidas's internal APIs by parsing `<script>` data from the HTML. These APIs could be used to fetch clean JSON data — a big win!

- **Rate Limiting & Bans**:  
  Even with pauses, repeated scraping attempts eventually led to 403 responses and temporary bans, making it difficult to consistently gather product data.

---

## ✅ Current Status

The data handling and parsing logic works well, but scraping is unstable due to anti-bot defenses. I ultimately submitted the project in a semi-working state after trying various bypass techniques.

---

## 💡 Notes

- Code is split into modular packages (`crawler`, `scraper`, `model`, etc.)
- Chromedp was the most effective tool but requires careful request pacing
- Output is saved to Excel using the [`excelize`](https://github.com/xuri/excelize) Go package
- Project is meant more as a learning exploration than a robust production scraper

---

## 🧠 Future Improvements

- Rotate proxies / user agents
- Use stealth plugins or headless Chrome patches
- Queue-based crawler architecture
- CAPTCHA bypass or human-in-the-loop fallback

---

Feel free to fork or learn from the code, and share your own scraping tips if you’ve cracked sites like Adidas successfully!
