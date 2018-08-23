angular.module('xlspaceship')
    .controller('XLSpaceshipCtrl', function($scope, $state, $http, $interval) {
        $scope.PLAYERID = "";
        $scope.PLAYERNAME = "";
        $scope.games = {};

        $scope.newOpponent = {
            host: "localhost",
            port: "8090",
        };

        function whoami() {
            $http.get("/xl-spaceship/user").catch(function(err) {
                console.log(err);
                alert(err.data || err);

                throw err
            }).then(function(res) {
                $scope.PLAYERID = res.data.user_id;
                $scope.PLAYERNAME = res.data.full_name;
            });
        };

        function challange() {
            $http.post("/xl-spaceship/user/game/new", {
                spaceship_protocol: {
                    hostname: $scope.newOpponent.host,
                    port: parseInt($scope.newOpponent.port, 10),
                }
            }, {headers: {'Content-Type': 'application/json'}}).catch(function(err) {
                console.log(err);
                alert(err.data || err);

                throw err
            }).then(function(res) {
                console.log(res.data);

                $scope.games[res.data.game_id] = res.data;

                $state.go('app.xlspaceship.play', {gameID: res.data.game_id});
            });
        }

        function refreshGame(gameID) {
            $scope.refreshing = true;

            return $http.get("/xl-spaceship/user/game/" + gameID).catch(function(err) {
                console.log(err);

                throw err;
            }).then(function(res) {
                console.log(res.data);

                return res.data;
            });
        }

        $scope.challange = challange;
        $scope.refreshGame = refreshGame;

        whoami();

        // @TODO: clear interval on $scope.$destroy
        $interval(function() {
            angular.forEach($scope.games, function(game, gameID) {
                refreshGame(gameID).then(function(game) {
                    $scope.games[gameID] = game;
                });
            });
        }, 1000);
    });
