Connect 4
--

This is an example of a simple turn-based, two player game to help show the basics of the game engine. The rules are simple and generally well known.

The rules for this game are entirely impleneted in a single custom stage (`stage.go`). This shows how player and server events are used to control the flow. Since this game has only one stage (not counting the JoinGame stage), it simply returns `nil` when the game is won to end the game event loop. A more complex game could use multiple game stages to build a kind of finite state machine in which case it would return the next stage instead of `nil`.