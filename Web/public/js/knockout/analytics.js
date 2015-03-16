/* GRAPH */
var helpedCountChart = new HelpedCountChart();
var helpedTimeChart = new HelpedTimeChart();
var peakChart = new PeakChart();


function AnalyticsViewModel() {
	var self = this;

	self.employees = [];
	self.helpCount = [];
	self.helpTime = [];
	self.peakTimes = [];
	
	/* 
	 *	Pull analytical data between given times
	 */
	self.pullData = function(ti, tf) {
		self.employees = [];
		self.helpCount = [];
		self.helpTime = [];
		self.peakTimes = [];
		self.REST = 0;
		self.REST_total = 3;

		$.post("/analytics/customersHourly",
			{"ti":ti, "tf":tf},
			function(data) {
				$.each(JSON.parse(data), function(key, value) {
				 	self.peakTimes.push([key,value])
				})
				peakChart.drawChart("#peakHours svg",
					peakChart.formatData(self.peakTimes))
			}
		);
		$.post("/analytics/helpCount",
    		{"ti":ti, "tf":tf},
    		function(data) {
				$.each(JSON.parse(data), function(key, value) {
					self.helpCount.push([key,value])
				})
				self.REST++;
				self.mapData();
    		}
		);
		$.post("/analytics/helpTime",
			{"ti":ti, "tf":tf},
			function(data) {
				$.each(JSON.parse(data), function(key, value) {
					self.helpTime.push([key,value])
				})
				self.REST++;
				self.mapData();
			}
		);
		$.get("/employees", function(data) {
				$.each(JSON.parse(data), function(i, employee) {
					self.employees.push(new Employee(employee));
				})
				self.REST++;
				self.mapData();
			}
		);
	};

	/*
	 *	Map analytical data to employee
	 */
	self.mapData = function() {
		if (self.REST < self.REST_total) return;

		/** Map data **/
		$.each(self.helpTime, function(i, data) {
			$.map($.grep(self.employees, function(e) {
				return e._id.$oid == data[0];
			}), function( val, ind) {
				val.helpTimes = data[1];
			});
		});
		
		
		$.each(self.helpCount, function(i, data) {
			$.map($.grep(self.employees, function(e) {
				return e._id.$oid == data[0];
			}), function( val, ind) {
				val.helpCount = data[1];
			});
		});

		/** Update charts **/
		helpedCountChart.drawChart("#helpCount svg", 
			helpedCountChart.formatData(self.employees));
		helpedTimeChart.drawChart("#helpTime svg",
			helpedTimeChart.formatData(self.employees));
	};
}

var vm = new AnalyticsViewModel();
$(function() {
	ko.applyBindings(vm);
});

var formatDate = function(d) {
	return String.leftPad(d.getMonth() + 1, 2, '0') + '-' + String.leftPad(d.getDate(), 2, '0') + '-' + d.getFullYear();
}

/* Date time picker */
$("#datetimepicker").datetimepicker({
	format:'m-d-Y',
	lang:'en',
	timepicker:false
});
$("#datetimepicker").val(formatDate(new Date()));
$("#datetimeselected").click(function() {
	m = (new Date($("#datetimepicker").val())).getTime()
	vm.pullData(m, m+24*3600*1000);
});
$("#datetimeselected").click();