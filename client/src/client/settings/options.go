package settings

type Options struct {
	Dir      string `short:"d" long:"dir" default:"data/" description:"Directory to watch and sync"`
	ClientId string `short:"c" long:"cid" description:"Client Id unique to this machine"`
	UserId   string `short:"u" long:"uid" description:"User Id unique to the user running using the client"`
}
