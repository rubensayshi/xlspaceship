XL Spaceship
============



Backend: Go
-----------
I decided to use Go as the language for the "backend" process, the main reason was just; because I can.

Normally Go wouldn't be my goto choice for a desktop game, currently the available cross-platform GUI libraries (Qt, etc)
are all early stage and experimental and most of them are a hassle to use and extremely verbose.
So a language with good and easy to use libraries for a cross-platform GUI would be preferred.

However nowadays more and more desktop applications actually build their GUI using HTML/CSS/JS
and the client simply bundles a lightweight browser to run the GUI (eg; https://github.com/asticode/go-astilectron).
With a javascript frontend development tends to be much quicker and the pool of developers available to work on it is also much larger.

Considering I haven't had any experience with any of these new-age GUI libs I went for a slightly simpler step
and just pop open a browser window with the GUI, moving forward it shouldn't be more than a few hours of fiddling with a library like `go-astilectron`
to turn it into more of a real desktop app.

Limited by the spec of the game it was a requirement that the background process handles of lot of the logic,
honestly if given free reign and in a real world scenario where a lot of the game logic would be managed by a server instead of a local instance
I'd simply write the whole client in javascript.

Go would be a good language for a game server backend to handle 1000s of players playing because
it's a good balance between performance (C-like languages) and productivty (scripting language),
so if this was a PoC for a real game it wouldn't be a bad choice to already write this in Go so the code can later be used in the actual server code.


#### The Project
The code has been split into 2 parts: `pkg/ssgame` and `pkg/ssclient`.


##### `pkg/ssgame`
The `pkg/ssgame` part contains all the code for 1 single game.

Currently it contains some code to use deterministic random values because that makes debugging and testing a lot easier,
for a real game it would ofcourse be important that the random numbers are secure, for which we'd need to switch to use the `crypto/rand` package.

Looking over it now I think it would have been better to simply implement a grid with coords and "place" the ships, hits and misses on them instead of tracking them seperately,
I haven't gotten around to doing that but with the test coverage in place it shouldn't be too hard to refactor that without breaking anything.


#### `pkg/ssclient`
The `pkg/ssclient` part contains the code to create and play games between players.

The HTTP request serving / handling in `server.go` is rather verbose atm, it could really use a library / framework to make the code less verbose with less boilerplate code needed.

Currently the `server.go` code does not have test coverage because without a library / framework it's a lot of hassle to write tests for the HTTP requests / responses,
it should absolutely be added, but considering it's 99% boilerplate code I've chosen not due to time constraints.

I've detached the `xlspacehip.go` from `server.go` using a channel to pass requests and responses back and forth between the 2,
combined with abstracting away the `requester` this allows us to swap out the `HTTPRequester` with a `MemRequester`
so that we can let 2 instaces of `XLSpaceship` communicate with each other as if they are sending HTTP requests, this is extremely useful for writing tests (see `xlspaceship_fullgame_new_test.go`).


#### Bundling statics into the binary
To bundle the static files into the binary that we spread to users we're using the `statik` package,
this package will go through our static files and place them all into a big blob in `./statik/statik.go` which will be used to serve the files from.

To make our live easy during development so we don't have to rebuild the go code every time we make changes to our GUI javascript code it will by default just serve the files from the filesystem directly.

We're using the build tag `statik` to switch between `server_fs.go` and `server_statik.go` when compiling to make it easy to switch between these 2 and to completely exclude the other method from the resulting binary.



Frontend: AngularJS 1.4
-----------------------
I decided to use AngularJS 1.4 for the frontend, as already described above using javascript for the GUI has some nice advantages,
on top of those mentioned I also had a small boilerplate project to use to save me time and it's very easy to read for reviewers as well.

For a real world project the choice would most likely be AngularJS 2 or ReactJS since those 2 are the current flavor of the month in the javascript world
and more importantly they're MUCH easier to write tests for and provide better performance (which wouldn't be relevant for such a small game but would be with a bigger game).
Alternatively there's WebGL for a mroe graphics heavy game (which has Go bindings as well with GopherJS).

I've decided to keep the GUI code rather simple due to time constraints and job I'm applying for is backend focused,
the GUI code also lacks tests, largely because of time constraints, to make the code easier to test a lot of it should be moved to services.

It would probably also be better to use a websocket to communicate with the backend process instead of the REST endpoints to avoid constant polling of the backend and have quicker response times,
but then the REST endpoints that were in the spec wouldn't be used anymore ...
however considering we've already detached the HTTP handle in the backend from the actual logic it shouldn't be much effort to do so.

All the code for the gui is in `./gui`.

#### Gulp
It uses `gulp` to build the files in `./gui/src`, again webpack is more flavor of the month but gulp is easy to read and understand and I had a boilerplate config for it ready to use.
After building the final files are in `./gui/web` and ready to be served by a HTTP file server.

#### The Project
To review all the relevant logic take a look at `./gui/src/js/controllers/` and `./gui/src/js/templates/xlspaceship`.



Future Improvements
-------------------
 - `/user` endpoint should have some way to authenticate that requests aren't coming from a different source than the GUI.
   this could easily be done by letting the background process inject an authentication token into the HTML before the page is opened.

 - `/protocol` should have some way to authenticate that requests are coming from the player of the game it concerns.
   this could easily be done by sharing an authentication token during the initial create game request.

 - without a central server using the user_id is problematic since both users can have the same user_id.
   currently it will reject creating a game when the user_id of the opponent is the same as your own.

 - without a central server it's very hard to know if your opponent isn't cheating,
   we could share a hash of our board (with a salt to avoid rainbow tables) when the game is created to verify the result at the end of a game.
   however, considering how few possibilities of spaceship layouts there are,
   the hash operation would require enough itterations to make sure you can't bruteforce your opponents board within the time of 1 game.