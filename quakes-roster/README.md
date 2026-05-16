quakes-roster generates a roster for the SJ Earthquakes using the info on their site.

run this to fetch information
```
go run . scrape
```

This downloads player info into `player.json` as well as images.

You may want to massage `player.json` a bit to give it the order and content you want.

Then run this to generate a pdf (`roster.pdf`):
```
go run . pdf
```
