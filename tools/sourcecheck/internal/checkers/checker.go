package checkers

type Checker interface {
	Check(fileName string) Messages
	Setup(dirName string, packageName string) error
}
