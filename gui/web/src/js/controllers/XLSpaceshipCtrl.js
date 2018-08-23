angular.module('xlspaceship')
    .controller('XLSpaceshipCtrl', function($scope, $state, $http) {
        $scope.PLAYERID = "";
        $scope.PLAYERNAME = "";
        $scope.games = {};

        $scope.newOpponent = {
            host: "localhost",
            port: "8090",
        };

        $scope.whoami = function() {
            $http.get("/xl-spaceship/user").catch(function(err) {
                console.log(err);
                alert(err.data || err);

                throw err
            }).then(function(res) {
                $scope.PLAYERID = res.data.user_id;
                $scope.PLAYERNAME = res.data.full_name;
            });
        };

        $scope.whoami();

        $scope.challange = function() {
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
        };
    });
