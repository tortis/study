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


StudyApp.directive('fileModel', ['$parse', function ($parse) {
	return {
		restrict: 'A',
		link: function(scope, element, attrs) {
			var model = $parse(attrs.fileModel);
			var modelSetter = model.assign;

			element.bind('change', function() {
				scope.$apply(function() {
					modelSetter(scope, element[0].files[0]);
				});
			});
		}
	};
}]);

StudyApp.controller('DeckCtrl', ['$scope', '$http', '$routeParams',
    function($scope, $http, $routeParams) {
        console.log("DeckCtrl started");
		$scope.newc = null;
		$scope.newimg = null;
        $scope.did = $routeParams.did;

		$scope.update = function() {
        	$http.get("/decks/"+$scope.did).
        	success(function(data,status, headers, config) {
        	    $scope.deck = data;
                for (var i = 0; i < $scope.deck.cards.length; i++) {
                    $scope.deck.cards[i].notess = $scope.deck.cards[i].notes.split("\r\n");
                }
        	}).
        	error(function(data, status, headers, config) {
        	    console.log("Error on call to /decks/test2");
        	});
		};

		$scope.update();

		$scope.postCard = function() {
			var fd = new FormData();
			fd.append('image', $scope.newimg);
			fd.append('json', JSON.stringify($scope.newc));
			console.log("Posting card: " + JSON.stringify($scope.newc));
			$http.post("/decks/"+$scope.did+"/cards", fd, {
				transformRequest: angular.identity,
				headers: {'Content-Type': undefined}
			})
			.success(function() {
				console.log("New card posted successfully")
				$scope.update();
				$scope.newc = null;
				$scope.newimg = null;
			})
			.error(function() {
				console.log("Failed to post new card.");
			});
		};

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
		$scope.updateCard = function(c, i) {
			$http.put("/decks/"+$scope.did+"/cards/"+i, JSON.stringify(c))
			.success(function(data, status, headers, config) {
				c.editing = false;
				console.log("Card updated successfully.");
			})
			.error(function(data, status, headers, config) {
				console.log("Failed to update card.");
			});	
		};
    }]);
