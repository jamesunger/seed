'use strict';

angular.module('seedApp')
  .controller('UserCtrl', function ($scope, $resource, $routeParams) {
	if ($routeParams.field && $routeParams.value) {
	     var Users = $resource('/user', {field: $routeParams.field, value: $routeParams.value},
		{ 'get': { method: 'GET', isArray: true}})
		var users = Users.get(function() {
			$scope.users = users;
		});
	} else {
	     var Users = $resource('/user', {},
		{ 'get': { method: 'GET', isArray: true}})
	     var users = Users.get(function() {
			$scope.users = users;
	        });
	}
  });
