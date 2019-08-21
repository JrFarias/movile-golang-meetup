func (c *consumer) Consume() {
  for w := 1; w <= c.workerPool; w++ {
   go c.worker(w)
  }
 }

 func (c *consumer) worker(id int) {
  for {
   output, err := retrieveSQSMessages(c.QueueURL, maxMessages)
   if err != nil {
    continue
   }

   var wg sync.WaitGroup
   for _, message := range output.Messages {
    wg.Add(1)

    go func(m *message) {
      defer wg.Done()
      if err := h(m); err != nil {
        continue
      }

      c.delete(m) 
    }(newMessage(m))
    
    wg.Wait()
   }
  }
 }


/////////////////////////////
////////// METRICS //////////
///////// 6,750 mpm /////////
/////////////////////////////
