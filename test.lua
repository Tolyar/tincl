test = ReadFromTelnet()
print ("* SMTP Greeting: ", test)
WriteToTelnet("helo qq")
test = ReadFromTelnet()
print ("* SMTP Answer after helo: ", test)

