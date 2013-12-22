'use strict';

angular.module('seedApp')
  .controller('UserCtrl', function ($scope, $resource) {
	//alert($resource)
	var Users = $resource('/user');
	var user = Users.get(function() {
		//alert(user);
		//alert(user.Address);
		//return user;
		$scope.user = user;
		//alert(user.Username);
			
	});
	//$scope.users = user;
	//alert(user.Username);
	//alert($scope.users.Username);
	//console.log($scope.users.Username);
  });
