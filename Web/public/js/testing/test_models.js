QUnit.test( "Employee", function( assert ) {
	var data = {_id: "employee", retailer: "UAlberta", name: "John", mac: "EOF"};
	var e = new Employee(data);

	assert.ok( e._id == data._id, "Passed!" );
	assert.ok( e.retailer == data.retailer, "Passed!" );
	assert.ok( e.name() == data.name, "Passed!" );
	assert.notOk( e.auth_code() == null, "Passed!" );
	assert.ok( e.mac() == data.mac, "Passed!" );
});

QUnit.test("Drawable", function( assert ) {
	var data = {mac: "drawable", x: 3.0, y: 4.0, radius:0.5, priority:0.5, employee:true};
	var d = new Drawable(data);

	assert.ok( d.id() == data.mac, "Passed!");
	assert.ok( d.x == data.x, "Passed!");
	assert.ok( d.y == data.y, "Passed!");
	assert.ok( d.radius == data.radius, "Passed!");
	assert.ok( d.priority == data.priority, "Passed!");
	assert.ok( d.employee == data.employee, "Passed!");
	assert.notOk( d.lastUpdated == null, "Passed!")
	assert.notOk( d.name() == null, "Passed!")
	assert.notOk( d.color() == null, "Passed!")

	var data2 = {x: 5.0, y: 1.1, radius:0.2, priority:0.7, lastUpdated:1000};
	d.update(data2);

	assert.ok( d.x == data2.x, "Passed!");
	assert.ok( d.y == data2.y, "Passed!");
	assert.ok( d.radius == data2.radius, "Passed!");
	assert.ok( d.priority == data2.priority, "Passed!");
	assert.ok( d.lastUpdated == data2.lastUpdated, "Passed!");
	assert.ok( d.employee == data.employee, "Passed!");
});

QUnit.test("Generate ID", function( assert ) {
	var g = generateID();

	assert.ok(typeof(g) == "string", "Passed!");
	assert.ok(g.length > 0, "Passed!");
});