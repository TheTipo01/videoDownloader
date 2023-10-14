package main

type config struct {
	Token      string   `fig:"token" validate:"required"`
	LogLevel   string   `fig:"loglevel"`
	URLs       []string `fig:"urls" validate:"required"`
	Channel    int64    `fig:"channel" validate:"required"`
	Downloader string   `fig:"downloader"`
}

type Downloader struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Endpoint  string `json:"endpoint"`
	Url       string `json:"url"`
	Type      string `json:"type"`
	VideoData struct {
		WmVideoUrl    string `json:"wm_video_url"`
		WmVideoUrlHQ  string `json:"wm_video_url_HQ"`
		NwmVideoUrl   string `json:"nwm_video_url"`
		NwmVideoUrlHQ string `json:"nwm_video_url_HQ"`
	} `json:"video_data"`

	ImageData struct {
		NoWatermarkImageList []string `json:"no_watermark_image_list"`
		WatermarkImageList   []string `json:"watermark_image_list"`
	} `json:"image_data"`
}
