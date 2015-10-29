# goreaders

Collection of chainables Go Readers with Cursor/Offset/Retry capabilities.
Include :
 - File Reader
 - Gzip Reader
 - Recursive Directory Reader
 - JSon Reader
 
Advantages:
 - Errors are produced only when reading or getting the cursor.
 - Recursive Directory Reader automatically manage gzip files
 - Simple cursor (string)
 - Offset and Start functions enable to retry units or batch of data
