package main

import (
	"flag"
	"fmt"
	"go-daemon/configure"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var h, s, i, u bool
	flag.BoolVar(&h, "h", false, "show help")
	flag.BoolVar(&s, "s", false, "run as service")
	flag.BoolVar(&i, "i", false, "install service")
	flag.BoolVar(&u, "u", false, "uninstall service")
	flag.Parse()

	if h {
		flag.PrintDefaults()
	}
	if e := configure.Init(); e != nil {
		log.Fatalln(e)
	}

	if e := SetToken(); e != nil {
		log.Println(e)
	}

	if i {
		install()
	} else if u {
		unstall()
	} else if s {
		runService()
	} else {
		runNormal()
	}

}
func install() {
	fmt.Println("/****	install	****/")
	cnf := configure.GetConfigure()
	fmt.Println(cnf)

	bin, e := filepath.Abs(os.Args[0])
	if e != nil {
		log.Fatalln(e)
	}

	e = InstallService(cnf.Name,
		cnf.Show,
		cnf.Description,

		bin+" -s",

		cnf.Auto,
	)
	if e == ErrorChangeServiceConfig2 {
		log.Println("install success")
		log.Println(e)
	} else if e != nil {
		log.Fatalln(e)
	}
}
func unstall() {
	fmt.Println("/****	unstall	****/")
	cnf := configure.GetConfigure()
	fmt.Println(cnf)

	e := UnstallService(cnf.Name)
	if e != nil {
		log.Fatalln(e)
	}
}
func runService() {
	fmt.Println("/****	runService	****/")
	cnf := configure.GetConfigure()
	fmt.Println(cnf)

	e := InitWorkDirectory()
	if e != nil {
		log.Fatalln(e)
	}

	s := newService()
	service_main(s, cnf.Name)
}
func runNormal() {
	fmt.Println("/****	runNormal	****/")
	cnf := configure.GetConfigure()
	fmt.Println(cnf)

	e := InitWorkDirectory()
	if e != nil {
		log.Fatalln(e)
	}

	s := newService()
	e = s.Run()
	if e != nil {
		log.Fatalln(e)
	}

	s.Wait()
}
