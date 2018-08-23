angular.module('xlspaceship')
    .controller('XLSpaceshipPlayCtrl', function($scope, $stateParams, $http, $timeout) {
        $scope.refreshing = false;
        $scope.game = $scope.games[$stateParams.gameID];

        $scope.salvo = randomSalvo();

        function randomShot() {
            let x = randomIntFromInterval(0, 15);
            let y = randomIntFromInterval(0, 15);

            return x.toString(16) + "x" + y.toString(16);
        }

        function randomSalvo() {
            return {
                shot1: randomShot(),
                shot2: randomShot(),
                shot3: randomShot(),
                shot4: randomShot(),
                shot5: randomShot(),
            };
        }

        function refresh() {
            $scope.refreshing = true;

            return $scope.refreshGame($stateParams.gameID).then(function(game) {
                $scope.game = game;

                $timeout(function() {
                    $scope.refreshing = false;
                }, 200);

                return game;
            });
        }

        function fireSalvo() {
            $http.put("/xl-spaceship/user/game/" + $stateParams.gameID + "/fire", {
                salvo: [$scope.salvo.shot1, $scope.salvo.shot2, $scope.salvo.shot3, $scope.salvo.shot4, $scope.salvo.shot5, ],
            }, {headers: {'Content-Type': 'application/json'}}).catch(function(err) {
                console.log(err);
                alert(err.data || err);
            }).then(function(res) {
                console.log(res.data);

                // new random salvo for next round
                $scope.salvo = randomSalvo();

                return refresh();
            });
        }

        $scope.refresh = refresh;
        $scope.fireSalvo = fireSalvo;

        if (!$scope.game) {
            refresh()
                .then(function(game) {
                    $scope.games[game.game_id] = game;
                }, function() {
                    $state.go('app.xlspaceship.welcome');
                })
            ;
        }
    });

