// we're gonnaa programmatically interact with Three.js scene.
	// adding objects and spraying magic dust
	scene = document.querySelector("a-scene")
	
	// scatter plot is a parent 3D object with nested sphere points.
	// you can move the scatter_plot object, and all the points will follow.
	scatter_plot = new THREE.Object3D()
	scene.object3D.add(scatter_plot)
	function create_sphere_point(radius) {
		s = new THREE.Mesh(
			new THREE.SphereGeometry(radius || 1, 16, 16),
			new THREE.MeshBasicMaterial({color: 0xffffff},
		))
		scatter_plot.add(s)
		return s
	}
	scatter_plot_points = []
	for (i=0; i < 100; i++) {
		for (j=0; j < 10; j++) {
			sphere = create_sphere_point(0.1)
			scatter_plot_points.push(sphere)
			sphere.position.z = 15 - i * .3
			sphere.position.x = j - 4.5
		}
	}
	scatter_plot.position.z = -20
	
	  
	// alternatively we tried to work with a mesh and manipulate it's points.
	// spoiler alert: the vertex coloring part was a mess and didn't work out in one night.
	// spectrograph_plane = document.getElementById("spectrograph_plane")
	// spectrograph_geometry = spectrograph_plane.object3DMap.mesh.geometry


        // provides a stream of data after main.py has been run.
	var ws = new WebSocket("ws://alpine-unicorn.noise/ws")


	// data coming in from ECG will be stored in readings_history ledger.
	// every reading is an array of ten values
	// we'll try to keep only a 100 of records
	// POOR CODE QUALITY WARNING: this magic number 100 is also hardcoded into scatter plot generation
	// so... just beware
	readings_history = []
	function get_reading_data(gen) {
		return readings_history[gen] || [0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
	}

	// Color Lookup Table is an awesome thing
	// you create one with one of awailable spectrums, define the amount of
	// color gradations and a linear scale to map to.
	// you can then get a color for a specific value like Lut.getColor(15.67)
	colors_lut = new THREE.Lut('rainbow', 512)
	colors_lut.setMin(0)
	colors_lut.setMax(6.5)

	ws.onmessage = function (event) {
	
		// created an array here to read out data from the arduino chip.
		// We split with a space. The first value is the right reader, 
		// the second is the left reader.
		var data = JSON.parse(event.data);
		if (data.Name == "fft") {
			var fft_readings = data.Channels;
			// Normalize the channels
			// Ex: [1, 2, 3, 4] => [0.1, 0.2, 0.3, 0.4]
			// Ex: [10, 20, 30, 40] => [0.1, 0.2, 0.3, 0.4]
			var lchan_sum = fft_readings[0].reduce(function(x, y) { return x + y }, 0)
			fft_readings[0] = fft_readings[0].map(function(x) { return x / lchan_sum })
			var rchan_sum = fft_readings[1].reduce(function(x, y) { return x + y }, 0)
			fft_readings[1] = fft_readings[1].map(function(x) { return x / rchan_sum })
			var lchans = {"gamma":fft_readings[0].slice(30, 45, fft_readings[1].length).reduce(function(x, y) { return  x + y }, 0),
						   "beta" :fft_readings[0].slice(13, 30).reduce(function(x, y) { return x + y }, 0),
						   "alpha":fft_readings[0].slice(8, 12).reduce(function(x, y) { return x + y }, 0),
						   "theta":fft_readings[0].slice(4, 7.75).reduce(function(x, y) { return x + y }, 0),
						   "delta":fft_readings[0].slice(2, 4).reduce(function(x, y) { return x + y }, 0)}
			var rchans = {"gamma":fft_readings[1].slice(30, 45, fft_readings[1].length).reduce(function(x, y) { return x + y }, 0),
						   "beta" :fft_readings[1].slice(13, 30).reduce(function(x, y) { return x + y }, 0),
						   "alpha":fft_readings[1].slice(8, 12).reduce(function(x, y) { return x + y }, 0),
						   "theta":fft_readings[1].slice(4, 7.75).reduce(function(x, y) { return x + y }, 0),
						   "delta":fft_readings[1].slice(2, 4).reduce(function(x, y) { return x + y }, 0)}
		} else if (data.Name == "sample") {
			var eeg_readings = data.Channels;
			return
		}

		// changing the string to a decimal by dividing or 
		// multiplying the strings given by our raw data
		// there is too much variance. Division, modifies the range of physical 
		// bounds variance but not the time update variance
		var lgamma = lchans.gamma
		var rgamma = rchans.gamma		
		var lbeta = lchans.beta
		var rbeta = rchans.beta
		var lalpha = lchans.alpha
		var ralpha = rchans.alpha
		var ltheta = lchans.theta
		var rtheta = rchans.theta
		var ldelta = lchans.delta
		var rdelta = rchans.delta
		
		// preparing and storing the readings data
		data = [rgamma, rbeta, ralpha, rtheta, rdelta, ldelta, ltheta, lalpha, lbeta, lgamma]
		readings_history.unshift(data)
		readings_history = readings_history.slice(0, 100)

		// animating methes, points, etc.
		for (generation = 0; generation < 100; generation++) {
			data = get_reading_data(generation)
		        for (idx in data) {
			  // calculate the vertex/point index that we're going to update,
			  // as well as the color
			  vertice_idx = generation * 10 + (+idx)
			  scaled_value = data[idx] * 10
			  point_color = colors_lut.getColor(scaled_value)

			  // here we manipulated the spectrograph mesh geometry and vertex colors.
			  // the latter one didn't quite work out
			  /* spectrograph_geometry.vertices[vertice_idx].z = scaled_value
			  //if (scaled_value > 0) console.log(scaled_value, point_color)
			  face = spectrograph_geometry.faces[generation * 9 + 2 * (+idx)]
			  face.color.set(point_color)
			  face = spectrograph_geometry.faces[generation * 9 + 2 * (+idx) + 1]
			  face.color.set(point_color) */

			  // updating point cloud
			  scatter_plot_points[vertice_idx].position.y = scaled_value
			  scatter_plot_points[vertice_idx].material.color = point_color
			}
		}
	
		// this is required to notify mesh that it's internal geometry was updated. 
		// spectrograph_geometry.verticesNeedUpdate = true
		// spectrograph_geometry.colorsNeedUpdate = true
		
		// scatter_plot: from Russia with love
		// https://github.com/vintlucky777
		// vintlucky777@gmail.com
	};
