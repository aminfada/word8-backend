package cmd

import (
	"fmt"
	"strconv"
	"time"
	configs "vocab8/config"
	"vocab8/vocab"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var cfgPath string

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run an instance of the service",
	Long:  "",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	fmt.Println("Initializing the Service...")

	configs.ParseConfig(cfgPath, &configs.Cfg)
	configs.NewDB(*configs.Cfg)

	go func() {
		for {
			vocab.RenewThePool()
			time.Sleep(time.Minute * 15)
		}
	}()

	r := gin.Default()
	r.Use(vocab.CorsMiddleware())
	r.StaticFile(configs.Cfg.Context+"speech", configs.Cfg.SpeechPath+"speech.mp3")
	baseRoute := r.Group(configs.Cfg.Context)
	{
		baseRoute.GET("vocab", vocab.DrawVocabOnMap)
		baseRoute.POST("vocab", vocab.AddVocab)
		baseRoute.PUT("vocab", vocab.SubmitFeedback)
	}

	fmt.Println("the service is live.")
	r.Run(configs.Cfg.Host + ":" + strconv.Itoa(configs.Cfg.Port))

}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVar(&cfgPath, "conf", confDir, "The path containing the configuration of wallet service")

}
