package uadmin

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/otiai10/copy"
	"github.com/uadmin/uadmin/interfaces"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type MigrateCommand struct {

}
func (c MigrateCommand) Proceed() {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &CommandRegistry{
		actions: make(map[string]interfaces.CommandInterface),
	}
	createCommand := new(CreateMigration)
	createCommand.opts = &CreateMigrationOptions{}

	commandRegistry.addAction("create", interfaces.CommandInterface(createCommand))
	upCommand := new(UpMigration)
	upCommand.opts = &UpMigrationOptions{}

	commandRegistry.addAction("up", interfaces.CommandInterface(upCommand))
	downCommand := new(DownMigration)
	downCommand.opts = &DownMigrationOptions{}

	commandRegistry.addAction("down", interfaces.CommandInterface(downCommand))
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.isRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return
	}
	commandRegistry.runAction(action)
}

func (c MigrateCommand) ParseArgs() {

}

func (c MigrateCommand) GetHelpText() string {
	return "Migrate your database"
}

var migrationTplPath = "internal/templates/migrations"
var re = regexp.MustCompile("[[:^ascii:]]")

func prepareMigrationName(message string) string {
	now := time.Now()
	sec := now.Unix()
	message = re.ReplaceAllLiteralString(message, "")
	if len(message) > 10 {
		message = message[:10]
	}
	return fmt.Sprintf("%s_%d", message, sec)
}

type CreateMigrationOptions struct {
	Message string `short:"m" required:"true" description:"Describe what is this migration for"`
	Blueprint string `short:"b" required:"true" description:"Blueprint you'd like to create migration for'"`
}

type CreateMigration struct {
	opts *CreateMigrationOptions
}

func (command CreateMigration) ParseArgs() {
	parser := flags.NewParser(command.opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}

func (command CreateMigration) Proceed() {
	bluePrintPath := "blueprint/" + strings.ToLower(command.opts.Blueprint)
	if _, err := os.Stat(bluePrintPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Blueprint %s doesn't exist", command.opts.Blueprint))
	}
	dirPath := "blueprint/" + strings.ToLower(command.opts.Blueprint) + "/migrations"
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			panic(err)
		}
	}
	pathToBaseMigrationsFile := dirPath + "/migrations.go"
	if _, err := os.Stat(pathToBaseMigrationsFile); os.IsNotExist(err) {
		err := copy.Copy(migrationTplPath+"/migrations.go.tpl", pathToBaseMigrationsFile)
		if err != nil {
			panic(err)
		}
	}
	migrationName := prepareMigrationName(command.opts.Message)
	pathToConcreteMigrationsFile := dirPath + "/" + migrationName + ".go"
	var lastMigrationId int
	var err error
	err = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		var migrationFileRegex = regexp.MustCompile(`.*?_(\d+)\.go`)
		match := migrationFileRegex.FindStringSubmatch(path)
		if len(match) > 0 {
			migrationId, _ := strconv.Atoi(match[1])
			if migrationId > lastMigrationId {
				lastMigrationId = migrationId
			}
		}
		return nil
	})
	if _, err = os.Stat(pathToConcreteMigrationsFile); os.IsNotExist(err) {
		err := copy.Copy(migrationTplPath+"/concrete_migration.go.tpl", pathToConcreteMigrationsFile)
		if err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Sprintf("migration %s already exists", pathToConcreteMigrationsFile))
	}
	migrationConcreteString, err := ioutil.ReadFile(migrationTplPath+"/concrete_migration.go.tpl")
	if err != nil {
		panic(err)
	}
	humanizedMessage := strings.ReplaceAll(
		re.ReplaceAllLiteralString(command.opts.Message, ""),
		"\"",
		"",
	)
	migrationConcreteStringNew := strings.Replace(
		string(migrationConcreteString),
		"concreteMigrationName",
		humanizedMessage, -1)
	now := time.Now()
	sec := now.Unix()
	migrationConcreteStringNew = strings.Replace(
		string(migrationConcreteStringNew),
		"concreteMigrationId",
		strconv.Itoa(int(sec)), -1)
	if lastMigrationId > 0 {
		migrationConcreteStringNew = strings.Replace(
			string(migrationConcreteStringNew),
			"dependencyId",
			fmt.Sprintf("[]string{\"%s\"}", strconv.Itoa(lastMigrationId)), -1)
	} else {
		migrationConcreteStringNew = strings.Replace(
			string(migrationConcreteStringNew),
			"dependencyId",
			"make([]string, 0)", -1)
	}
	migrationConcreteStringNew = strings.Replace(
		string(migrationConcreteStringNew),
		"MigrationName",
		migrationName, -1)
	err = ioutil.WriteFile(pathToConcreteMigrationsFile, []byte(migrationConcreteStringNew), 0755)
	if err != nil {
		panic(err)
	}
	migrationInitializationString, err := ioutil.ReadFile(migrationTplPath+"/migration_initialization.go.tpl")
	if err != nil {
		panic(err)
	}
	migrationInitializationStringNew := strings.Replace(
		string(migrationInitializationString),
		"migrationName",
		migrationName, -1)

	read, err := ioutil.ReadFile(pathToBaseMigrationsFile)
	if err != nil {
		panic(err)
	}
	newContents := strings.Replace(
		string(read),
		"// placeholder to insert next migration",
		migrationInitializationStringNew + "\n    // placeholder to insert next migration", -1)
	err = ioutil.WriteFile(pathToBaseMigrationsFile, []byte(newContents), 0755)
	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"Created migration for blueprint %s with name %s\n",
		command.opts.Blueprint,
		command.opts.Message,
	)
}

func (command CreateMigration) GetHelpText() string {
	return "Create migration for your blueprint"
}

type UpMigrationOptions struct {
}

type UpMigration struct {
	opts *UpMigrationOptions
}

func (command UpMigration) ParseArgs() {
	parser := flags.NewParser(command.opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}

func (command UpMigration) Proceed() {

}

func (command UpMigration) GetHelpText() string {
	return "Upgrade your database"
}

type DownMigrationOptions struct {
}

type DownMigration struct {
	opts *DownMigrationOptions
}

func (command DownMigration) ParseArgs() {
	parser := flags.NewParser(command.opts, flags.Default)
	_, err := parser.ParseArgs(os.Args[2:])
	if err != nil {
		log.Fatal(err)
	}
}

func (command DownMigration) Proceed() {

}

func (command DownMigration) GetHelpText() string {
	return "Downgrade your database"
}
