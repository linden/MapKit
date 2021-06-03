package mapkit

import (
	"github.com/linden/fetch"

	"fmt"
	"net/url"
)

const UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"

type MapKit struct {
	token   string
	refresh string
}

func (mapkit *MapKit) Refresh() error {
	request, err := fetch.Fetch("https://cdn.apple-mapkit.com/ma/bootstrap?apiVersion=2&mkjsVersion=5.61.1&poi=1", fetch.Options{
		Headers: fetch.Headers{
			"Connection":      "keep-alive",
			"Pragma":          "no-cache",
			"Cache-Control":   "no-cache",
			"User-Agent":      UserAgent,
			"Accept":          "*/*",
			"Origin":          "https://www.icloud.com",
			"Sec-Fetch-Site":  "same-site",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Dest":  "empty",
			"Referer":         "https://www.icloud.com/",
			"Accept-Language": "en-US,en;q=0.9",
			"authorization":   "Bearer " + mapkit.token,
		},
	})

	if err != nil {
		return err
	}

	if request.Status != 200 {
		return fmt.Errorf("invalid status: %s\n", request.Status)
	}

	body, err := request.JSON()

	if err != nil {
		return err
	}

	refresh, ok := body["accessKey"]

	if ok == false {
		return fmt.Errorf("invalid body: %s\n", body)
	}

	mapkit.refresh = refresh.(string)

	return nil
}

func (mapkit MapKit) GetTile(x int, y int, z int, scale int) (tile string, err error) {
	query := fmt.Sprintf("style=0&size=1&x=%d&y=%d&z=%d&scale=%d&lang=en&v=2105224&poi=0&accessKey=%s&labels=0&tint=light&emphasis=standard", x, y, z, scale, url.QueryEscape(mapkit.refresh))

	response, err := fetch.Fetch("https://cdn2.apple-mapkit.com/ti/tile?"+query, fetch.Options{
		Headers: fetch.Headers{
			"authority":       "cdn4.apple-mapkit.com",
			"pragma":          "no-cache",
			"cache-control":   "no-cache",
			"user-agent":      UserAgent,
			"accept":          "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
			"origin":          "https://www.icloud.com",
			"sec-fetch-site":  "same-site",
			"sec-fetch-mode":  "cors",
			"sec-fetch-dest":  "empty",
			"referer":         "https://www.icloud.com/",
			"accept-language": "en-US,en;q=0.9",
		},
	})

	if err != nil {
		return "", err
	}

	if response.Status != 200 {
		return "", fmt.Errorf("invalid status: %d\n", response.Status)
	}

	body, err := response.Text()

	if err != nil {
		return "", err
	}

	return body, nil
}

func (mapkit MapKit) GetSatelliteTile(x int, y int, z int, scale int) (tile string, err error) {
	query := fmt.Sprintf("style=7&size=2&x=%d&y=%d&z=%d&scale=%d&v=9082&poi=0&accessKey=%s", x, y, z, scale, url.QueryEscape(mapkit.refresh))

	response, err := fetch.Fetch("https://sat-cdn2.apple-mapkit.com/tile?"+query, fetch.Options{
		Headers: fetch.Headers{
			"authority":       "sat-cdn2.apple-mapkit.com",
			"pragma":          "no-cache",
			"cache-control":   "no-cache",
			"user-agent":      UserAgent,
			"accept":          "image/avif,image/webp,image/apng,image/svg+xml,image/*,*/*;q=0.8",
			"origin":          "https://www.icloud.com",
			"sec-fetch-site":  "same-site",
			"sec-fetch-mode":  "cors",
			"sec-fetch-dest":  "empty",
			"referer":         "https://www.icloud.com/",
			"accept-language": "en-US,en;q=0.9",
		},
	})

	if err != nil {
		return "", err
	}

	if response.Status != 200 {
		return "", fmt.Errorf("invalid status: %s\n", response.Status)
	}

	body, err := response.Text()

	if err != nil {
		return "", err
	}

	return body, nil
}

func Create(arguments ...string) (mapkit MapKit, err error) {
	if len(arguments) == 2 {
		return MapKit{arguments[0], arguments[1]}, nil
	}

	response, err := fetch.Fetch("https://setup.icloud.com/setup/ws/1/mapkitToken", fetch.Options{
		Headers: fetch.Headers{
			"Connection":      "keep-alive",
			"Pragma":          "no-cache",
			"Cache-Control":   "no-cache",
			"User-Agent":      UserAgent,
			"Accept":          "*/*",
			"Origin":          "https://www.icloud.com",
			"Sec-Fetch-Site":  "same-site",
			"Sec-Fetch-Mode":  "cors",
			"Sec-Fetch-Dest":  "empty",
			"Referer":         "https://www.icloud.com/",
			"Accept-Language": "en-US,en;q=0.9",
		},
	})

	if err != nil {
		return MapKit{}, err
	}

	if response.Status != 200 {
		return MapKit{}, fmt.Errorf("invalid status: %d\n", response.Status)
	}

	body, err := response.JSON()

	if err != nil {
		return MapKit{}, err
	}

	token, ok := body["jwt"]

	if ok == false {
		return MapKit{}, fmt.Errorf("invalid body: %s\n", body)
	}

	return MapKit{token.(string), ""}, nil
}
