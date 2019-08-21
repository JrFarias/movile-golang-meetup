func (c *consumer) Consume() {
  for {
   output, err := receiveSQSMessages(c.QueueURL) 
   if err != nil {
    //log error
    continue
   }
   for _, m := range output.Messages {
     c.run(newMessage(m))
     if err := h(m); err != nil {
       //log error
       continue
     }
     
     c.delete(m)
   }
  }
 }


/////////////////////////////
////////// METRICS //////////
///////// 160 mpm ///////////
/////////////////////////////