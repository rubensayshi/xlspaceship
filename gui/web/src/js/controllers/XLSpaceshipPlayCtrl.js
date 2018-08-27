angular.module('xlspaceship')
    .controller('XLSpaceshipPlayCtrl', function($scope, $state, $stateParams, $http, $timeout, $interval) {
        $scope.refreshing = false;
        $scope.game = $scope.games[$stateParams.gameID];

        /**
         * generate a random shot
         */
        function randomShot() {
            let x = randomIntFromInterval(0, 15);
            let y = randomIntFromInterval(0, 15);

            return x.toString(16) + "x" + y.toString(16);
        }

        /**
         * generate a fresh random salvo
         */
        function randomSalvo(nShots) {
            let salvo = [];
            for (let i = 0; i < nShots; i++) {
                salvo.push(randomShot());
            }

            return salvo;
        }

        /**
         * refresh the game status
         */
        function refresh() {
            $scope.refreshing = true;

            return $scope.refreshGame($stateParams.gameID).then(function(game) {
                $scope.games[game.game_id] = game;
                $scope.game = game;

                $timeout(function() {
                    $scope.refreshing = false;
                }, 200);

                return game;
            });
        }

        /**
         * fire a salvo to the other player
         */
        function fireSalvo() {
            $http.put("/xl-spaceship/user/game/" + $stateParams.gameID + "/fire", {
                salvo: $scope.salvo,
            }, {headers: {'Content-Type': 'application/json'}}).catch(function(err) {
                console.log(err);
                alert(err.data || err);
            }).then(function(res) {
                console.log(res.data);

                return refresh().then(function() {
                    // new random salvo for next round
                    $scope.salvo = randomSalvo($scope.game.self.shots);
                });
            });
        }

        $scope.refresh = refresh;
        $scope.fireSalvo = fireSalvo;

        // if we're missing the game data then attempt to refresh it, if it fails we goto welcome screen
        if (!$scope.game) {
            refresh()
                .then(function(game) {
                    console.log('refreshed');

                    // assign a random salvo to our input state
                    $scope.salvo = randomSalvo($scope.game.self.shots);
                }, function() {
                    $state.go('app.xlspaceship.welcome');
                })
            ;
        } else {
            // assign a random salvo to our input state
            $scope.salvo = randomSalvo($scope.game.self.shots);
        }

        // setup interval to fetch fresh data
        let refreshInterval = $interval(function() {
            refresh();
        }, 500);

        // clear interval when $scope is destroyed to avoid zombie polling
        $scope.$on("$destroy", function() {
            $interval.cancel(refreshInterval);
        })
    });

