angular.module('xlspaceship')
    .controller('XLSpaceshipCtrl', function($scope, $state, $http, $interval) {
        $scope.PLAYERID = "";
        $scope.PLAYERNAME = "";
        $scope.games = {};

        $scope.newOpponent = {
            host: "localhost",
            port: "8090",
        };

        /**
         * fetch status about self, name, ID and list of games
         */
        function whoami() {
            $http.get("/xl-spaceship/user")
                .then(function(res) {
                    $scope.PLAYERID = res.data.user_id;
                    $scope.PLAYERNAME = res.data.full_name;

                    angular.forEach(res.data.games, function(gameID) {
                        refreshGame(gameID).then(function(game) {
                            $scope.games[gameID] = game;
                        });
                    });
                }, function(err) {
                    console.log(err);
                    alert(err.data || err);

                    throw err
                });
        }

        /**
         * challange another player
         */
        function challange() {
            $http.post("/xl-spaceship/user/game/new", {
                spaceship_protocol: {
                    hostname: $scope.newOpponent.host,
                    port: parseInt($scope.newOpponent.port, 10),
                }
            }, {headers: {'Content-Type': 'application/json'}})
                .then(function(res) {
                    console.log(res.data);

                    $scope.games[res.data.game_id] = res.data;

                    $state.go('app.xlspaceship.play', {gameID: res.data.game_id});
                }, function(err) {
                    console.log(err);
                    alert(err.data || err);

                    throw err
                });
        }

        /**
         * refresh a game's data
         */
        function refreshGame(gameID) {
            $scope.refreshing = true;

            return $http.get("/xl-spaceship/user/game/" + gameID)
                .then(function(res) {
                    console.log(res.data.game_id + ": " + (res.data.game.won || res.data.game.player_turn));

                    return res.data;
                }, function(err) {
                    console.log(err);

                    throw err;
                });
        }

        $scope.challange = challange;
        $scope.refreshGame = refreshGame;

        // fetch self data straight away
        whoami();

        // setup interval to fetch fresh data
        // @TODO: clear interval on $scope.$destroy
        $interval(function() {
            whoami();

            angular.forEach($scope.games, function(game, gameID) {
                if (!game.won) {
                    refreshGame(gameID).then(function (game) {
                        $scope.games[gameID] = game;
                    });
                }
            });
        }, 500);
    });
