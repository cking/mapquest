# MapQuest Google Go Client  [![GoDoc](https://godoc.org/github.com/cking/mapquest?status.svg)](https://godoc.org/github.com/cking/mapquest) [![Go Report Card](https://goreportcard.com/badge/github.com/cking/mapquest)](https://goreportcard.com/report/github.com/cking/mapquest) [![Sourcegraph](https://sourcegraph.com/github.com/cking/mapquest/-/badge.svg)](https://sourcegraph.com/github.com/cking/mapquest?badge)

This is a very first draft of a Google Go client for
[MapQuest Open Data Map APIs and Web Services](http://developer.mapquest.com/web/products/open).

## Status

We just implemented a limited set of APIs: The Static Map API,
the Geocoding API, and the Nominatim API. Other APIs will be
implemented as needed (pull requests are welcome).

Consider this package beta. The API is not stable and the code probably
is not production quality yet. We use it in parts of our applications,
but its use is limited. Bugs will be fixed when found. If you find a
bug, report it or--even better--send a pull request.

## Testing

To run the tests, you need to add a file `ACCESS_KEY` to the packages root
directory. Paste you MapQuest access key there.

Notice that the MapQuest API seems to not like the access key URL-encoded.
So make sure you paste it unencoded. For example, a valid access key should
look like a bit like this:

    Fmjad|lufd281r2q,72=o5-9attor

(Do not use the key above. It's just an example. The key above will not
yield valid results. Get your own key instead.
[It's free!](http://developer.mapquest.com/web/products/open))

After you created the file, you can run tests as usual:

    $ go test


## Creating a client

To use the various APIs, you first need to create a client.

    client := mapquest.NewClient("<your-app-key>")


Now that you have a Client, you can use the APIs.

## Static Map API

Here's an example of how to use the MapQuest static map API:

    req := &mapquest.StaticMapRequest{
      Center: "11.54165,48.151313",
      Zoom:   9,
      Width:  500,
      Height: 300,
      Format: "png",
    }
    img, err := client.StaticMap().Map(req)
    if err != nil {
      panic(err)
    }

You now have an [`image.Image`](http://golang.org/pkg/image/#Image) at hand.
Further details can be found in the
[Open Static Map Service Developer's Guide](http://open.mapquestapi.com/staticmap/).

## Geocoding API

The [Geocoding API](http://open.mapquestapi.com/geocoding/) enables you
to take an address and get the associated latitude and longitude.

    res, err := client.Geocoding().SimpleAddress("1090 N Charlotte St, Lancaster", 0)
    if err != nil {
      panic(err)
    }

Further details can be found in the
[Open Geocoding Service Developer's Guide](http://open.mapquestapi.com/geocoding/).

## Nominatim API

The [Nominatim API](http://open.mapquestapi.com/nominatim/) is a search
interface that relies solely on the data contributed to
[OpenStreetMap](http://www.openstreetmap.org/). It does not require an App-Key.

    res, err := client.Nominatim().SimpleSearch("Unter den Linden 117, Berlin, DE", 1)
    if err != nil {
      panic(err)
    }

Further details can be found in the
[Nominatim Search Service Developer's Guide](http://open.mapquestapi.com/nominatim/)

# Contributors

* [Oliver Eilhard](https://github.com/olivere/) (original author)
* [Christopher König](https://github.com/cking) (new author)

# License

This code comes with a [MIT
license](https://github.com/olivere/mapquest/blob/master/LICENSE).
