# maingo - delicious logo service <img src="logo.png" width="75px"/>

Maingo is a simple way to fetch logos and concatente them into a collage.

This is currently deployed at `http://maingo.ml`

How to call the service?
===
Example Request

    GET http://maingo.ml/logo/canvas?q=redis,scala,python
  
Explanation

| Query Parameter  | Explanation | Sample Value
| ---------------- | ------------| ------------- |
| `q` - required  | A list of comma separated logo names. For available logos, see [here](https://github.com/moo-mou/logo) | ?q=scala
| `width` - optional | A single number specifying the width of each logo. | ?q=scala&width=72  |
| `noCache` - optional | Set to `true` to regenerate an image | ?q=scala&noCache=true  |

How to run locally?
===
See [here](https://github.com/moo-mou/maingo#built-with) for dependencies.

Once all the dependencies are setup, you can run the program `./go-reload fishbowl/*`.

FAQ
===
Q: Where are logos coming from?

A: The logos are manually curated and available on Github [here](https://github.com/moo-mou/logo).

<br>

Q: Why <i>fishbowl</i>?

A: The idea is the API server can be extended with many services and they all live in the same bowl!

Example Use Case
===
[ Built with Collage ](https://github.com/moo-mou/maingo#built-with)

Built with
===
<img src="built_with.png"/>
