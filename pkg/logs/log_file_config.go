package logs

type LogFileConfig struct {
    Path       string `json:"path" default:"./logs/app.log"`
    MaxSize    int    `json:"maxSize" default:"100"`
    MaxBackups int    `json:"maxBackups" default:"10"`
    MaxAge     int    `json:"maxAge" default:"30"`
    Encode     string `json:"encode"  default:"json"`
}
