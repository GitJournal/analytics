3. Add a go handler to save incoming data to clickhouse
   - http post - https://github.com/pierrec/lz4 maybe?
   - it would be nice to use protobuf
4. Create a dart client to send this data
5. Get location from IP

On the dart side, keep serializing each event generated to a protofbuf file of when the program was launched
also store the other metadata required?

That way it'll be easy to know how much data is being sent?
Other things to do - ?

https://clickhouse.tech/docs/en/faq/integration/json-import/

# Small tasks
* Create an Event class
