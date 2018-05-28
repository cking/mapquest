/*
Package mapquest enables access to the Open MapQuest APIs.
For further details, see http://open.mapquestapi.com/.

To get started, you need to create a client:

	client := mapquest.NewClient("<your-app-key>")

Now that you have a client, you can use the APIs.

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
*/
package mapquest
