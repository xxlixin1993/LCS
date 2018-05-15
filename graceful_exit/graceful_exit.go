package graceful_exit

import (
	"container/list"
	"errors"
	"strings"
)

var exitList *ExitList

type ExitInterface interface {
	// get module name
	GetModuleName() string

	// exit module function
	Stop() error
}

type ExitList struct {
	// exit list
	ll *list.List

	// exit module name
	module map[string]*list.Element
}


// Initialize exit list.
func InitExitList() {
	exitList = &ExitList{
		ll:     list.New(),
		module: make(map[string]*list.Element),
	}
}

// Get a exitList instance
func GetExitList() *ExitList {
	if exitList == nil {
		return nil
	}
	return exitList
}

// Inserts a new element exitInterface at the front of exit list.
func (el *ExitList) Pop(exitInterface ExitInterface) error {
	if el.module == nil {
		return errors.New("[Smoothly Exit] Pop: plz init ExitList first")
	}

	// Judge whether it exists or not
	moduleName := exitInterface.GetModuleName()
	if _, ok := el.module[moduleName]; ok {
		return errors.New("[Smoothly Exit] Pop: this module(" + moduleName + ") name is exist")
	}

	// Add value
	element := el.ll.PushFront(exitInterface)
	el.module[moduleName] = element

	return nil
}

// Exit smoothly
func (el *ExitList) Stop() error {
	length := el.ll.Len()
	if length == 0 {
		return nil
	}

	errInfo := make([]string, 0)
	for i := 0; i < length; i++ {
		element := el.ll.Front()
		exitElement := element.Value.(ExitInterface)

		if err := exitElement.Stop(); err != nil {
			errInfo = append(errInfo, "[Smoothly Exit]: Stop this module("+exitElement.GetModuleName()+")"+err.Error())
		}

		el.ll.Remove(element)
	}

	if len(errInfo) > 0 {
		return errors.New(strings.Join(errInfo, "\n"))
	}

	return nil
}
