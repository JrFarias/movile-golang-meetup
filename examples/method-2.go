func (c *consumer) Consume() {
  c.workerPool = 50
  maxMessages := 10
  jobs := make(chan *message)
  for w := 1; w <= c.workerPool; w++ {
    go c.worker(w, jobs)
  }
  for {
    output, err := retrieveSQSMessages(c.QueueURL, maxMessages)
    if err != nil {
     continue
    }
    for _, m := range output.Messages {
      jobs <- newMessage(m)
    }
  }
}
func (c *consumer) worker(id int, messages <-chan *message) {
 for m := range messages {
   if err := h(m); err != nil {
    continue
   
   }
   c.delete(m)
 }
}

/////////////////////////////
////////// METRICS //////////
///////// 2,700 mpm /////////
/////////////////////////////
