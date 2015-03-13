var StudyApp = angular.module('StudyApp', ['ngRoute']);

StudyApp.config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        $routeProvider.
        when('/', {
            templateUrl: '/partials/decks.html',
            controller: 'DecksCrtl'
        }).
        when('/:did', {
            templateUrl: '/partials/deck.html',
            controller: 'DeckCtrl'
        });
    }]);

StudyApp.controller('DeckCtrl', ['$scope', '$http', '$routeParams',
    function($scope, $http, $routeParams) {
        console.log("DeckCtrl started");
        $scope.did = $routeParams.did;
        $http.get("/decks/"+$scope.did).
        success(function(data,status, headers, config) {
            $scope.deck = data;
            console.log(JSON.stringify(data));
        }).
        error(function(data, status, headers, config) {
            console.log("Error on call to /decks/test2");
        });

        $scope.deleteCard = function(i) {
            $http.delete("/decks/"+$scope.did+"/cards/"+i)
            .success(function(data, status, headers, config) {
                console.log("Successfully deleted card");
                $scope.deck.cards.splice(i, 1);
            })
            .error(function(data, status, headers, config) {
                console.log("Failed to delete card");
            });
        };
    }]);
