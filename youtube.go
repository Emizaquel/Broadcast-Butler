package main

func youtube_authenticate() {

}

func youtube_init() {
	youtube := platform{
		Name:       "YouTube",
		Url:        "https://www.googleapis.com/auth/youtube",
		Operations: []func(){youtube_authenticate},
	}
	platforms = append(platforms, youtube)
}
