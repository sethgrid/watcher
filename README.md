#watcher

A silly toy example of structuring code for testability. This tests a file watcher without using interfaces nor mocks. 

This is achieved by disabling the polling in the test instance and pushing poll results onto a channel directly, instead of the poller.