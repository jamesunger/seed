'use strict';

angular.module('seedApp')
  .controller('UserCtrl', function ($scope, $resource) {
	var Users = $resource('/user', {},
		{ 'get': { method: 'GET', isArray: true}})
	var users = Users.get(function() {
		$scope.users = users;
			
	});
  });
