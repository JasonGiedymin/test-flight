package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "runtime"
    "sync"
)

func watchForEventsOn(channel ApiChannel) {
    for msg := range channel {
        Logger.Trace("DOCKER EVENT:", *msg)
    }
}

func watchContainerOn(channel ContainerChannel, wg *sync.WaitGroup) {
    for msg := range channel {
        runtime.Gosched()
        Logger.Console(msg)
    }

    wg.Done()
}
