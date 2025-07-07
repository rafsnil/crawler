package common

//// Append a URL to visited.txt if not already present
//func SaveURLToFile(url, fileName string) error {
//	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	if err != nil {
//		return err
//	}
//	defer f.Close()
//
//	_, err = f.WriteString(url + "\n")
//	if err != nil {
//		return err
//	}
//
//	fmt.Printf("✅ Saved URL: %s\n", url)
//	return nil
//}
//
//// Load all visited URLs into a map for deduplication
//func LoadURLsFromFile(fileName string) (map[string]struct{}, error) {
//	urls := make(map[string]struct{})
//
//	file, err := os.Open(fileName)
//	if err != nil {
//		// If file doesn't exist, just return an empty map
//		if os.IsNotExist(err) {
//			return urls, nil
//		}
//		return nil, err
//	}
//	defer file.Close()
//
//	scanner := bufio.NewScanner(file)
//	for scanner.Scan() {
//		line := scanner.Text()
//		urls[line] = struct{}{}
//	}
//
//	if err := scanner.Err(); err != nil {
//		return nil, err
//	}
//
//	fmt.Printf("📄 Loaded %d visited URLs\n", len(urls))
//	return urls, nil
//}
