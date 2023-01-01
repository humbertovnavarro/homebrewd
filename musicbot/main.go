package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var IsLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString

func main() {
	godotenv.Load()
	cmd := exec.Command("java", "-Dnogui=true", "-jar", os.Getenv("JAR_FILE"))
	defer cmd.Process.Kill()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := generateConfig()
	if err != nil {
		panic(err)
	}
	go func() {
		err = cmd.Run()
		if err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/api/playlists", func(ctx *gin.Context) {
		playlists, err := os.ReadDir("./Playlists")
		playlistStrings := make([]string, len(playlists))
		for i, f := range playlists {
			playlistStrings[i] = f.Name()
		}
		if err != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		ctx.AbortWithStatusJSON(http.StatusOK, playlistStrings)
	})

	type PlaylistPost struct {
		Name    string   `json:"name"`
		Entries []string `json:"entries"`
	}

	r.POST("/api/playlists", func(ctx *gin.Context) {
		playlist := &PlaylistPost{}
		err := ctx.BindJSON(playlist)
		if err != nil {
			ctx.AbortWithStatus(http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		if !IsLetter(playlist.Name) {
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}

		f, err := os.Create(fmt.Sprintf("./Playlists/%s.txt", playlist.Name))
		if err != nil {
			ctx.AbortWithStatus(http.StatusNotFound)
		}
		f.WriteString(strings.Join(playlist.Entries, "\n"))
		ctx.AbortWithStatusJSON(http.StatusOK, playlist)
	})

	r.Static("/playlists", "./Playlists")
	r.Run()
}

var configTemplateString = `token = %s
owner = %s
prefix = "%s"
game = "DEFAULT"
status = ONLINE
songinstatus=false
altprefix = "NONE"
success = "🎶"
warning = "💡"
error = "🚫"
loading = "⌚"
searching = "🔎"
help = help
npimages = false
stayinchannel = false
maxtime = 0
alonetimeuntilstop = 0
playlistsfolder = "Playlists"
updatealerts=true
lyrics.default = "A-Z Lyrics"
aliases {
	settings = [ status ]
	lyrics = []
	nowplaying = [ np, current ]
	play = []
	playlists = [ pls ]
	queue = [ list ]
	remove = [ delete ]
	scsearch = []
	search = [ ytsearch ]
	shuffle = []
	skip = [ voteskip ]
	prefix = [ setprefix ]
	setdj = []
	settc = []
	setvc = []
	forceremove = [ forcedelete, modremove, moddelete ]
	forceskip = [ modskip ]
	movetrack = [ move ]
	pause = []
	playnext = []
	repeat = []
	skipto = [ jumpto ]
	stop = []
	volume = [ vol ]
}
eval=false
`

func generateConfig() error {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	owner := os.Getenv("DISCORD_BOT_OWNER")
	prefix := os.Getenv("DISCORD_BOT_PREFIX")
	f, err := os.Create("config.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	f.WriteString(fmt.Sprintf(configTemplateString, token, owner, prefix))
	if err != nil {
		return err
	}
	return nil
}
