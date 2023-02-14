# Go experimentation

## content
### event loop
very simple and stupid event loop to deal with the async calls using go routines and channels
refer to the [Test file ](./src/eventloop/scheduler_test.go) to see how it works

### retry 
a simple [retry](./src/retry/retry_test.go) logic with exponential backoff strategy

### httpclient
a simple abstraction for [http client](./src/httpclient/caller_test.go) which will handle HTTP/1.1 `request` | `response`  
 