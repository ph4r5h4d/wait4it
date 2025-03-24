package banner

import "fmt"

func Banner() {
	fmt.Println(applicationName() + sponsor() + support())
}

func applicationName() string {
	return " _       __            _    __    __ __    _    __ \n" +
		"| |     / /  ____ _   (_)  / /_  / // /   (_)  / /_\n" +
		"| | /| / /  / __ `/  / /  / __/ / // /_  / /  / __/\n" +
		"| |/ |/ /  / /_/ /  / /  / /_  /__  __/ / /  / /_  \n" +
		"|__/|__/   \\__,_/  /_/   \\__/    /_/   /_/   \\__/  \n\n"
}

func sponsor() string {
	return "You can buy me a coffee or sponsor wait4it via: \nhttps://paypal.me/ph4r5h4d \n" +
		"or\nhttps://github.com/sponsors/ph4r5h4d\n\n"
}

func support() string {
	return "For NFR and issues go to: https://github.com/ph4r5h4d/wait4it/issues\n\n"
}
