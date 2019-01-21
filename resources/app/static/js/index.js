const inputPixelSize = 20
const inputWidth = 300
const inputHeight = 300
const outputPixelSize = 4
const colorMap = {
	'white':  '#ffffff',
	'black':  '#222222',
	'red':    '#ce4a4a',
	'yellow': '#eaaf41',
	'green':  '#48a56a',
	'blue':   '#6688c3',
	'purple': '#b25da6',
};
let activeColor = null;
let currentRequestVersion = 0;

let index = {
	init: function() {
		// asticode.loader.init();
        asticode.modaler.init();
        asticode.notifier.init();

        index.initInput();
        index.initPalette();

        document.addEventListener('astilectron-ready', function() {
        	index.requestNewImg();
        })
	},
	beginLoading: function() {
		const loader = document.getElementById('loader')
		const cancelButton = document.getElementById('cancel')
		loader.classList.add('active')
		cancelButton.classList.add('active')
	},
	finishLoading: function() {
		const loader = document.getElementById('loader')
		const cancelButton = document.getElementById('cancel')
		loader.classList.remove('active')
		cancelButton.classList.remove('active')
	},
	requestNewImg: function() {
		index.beginLoading();
		currentRequestVersion++;
		const inputData = index.getInputData()
		const request = { name: 'new', payload: { 
			Data: Array.from(inputData.data), 
			Width: inputData.width, 
			Height: inputData.height,
			Version: currentRequestVersion,
		}}
		astilectron.sendMessage(request, function(response) {
			console.log(response)
			if (response.payload.version === currentRequestVersion) {
				index.finishLoading();
			} 
			if (response.name === 'error') {
				asticode.notifier.error(response.payload);
                return
			}
			if (!response.payload) {
				return
			}
			const responsePayload = response.payload
			if (responsePayload.version !== currentRequestVersion) {
				return
			}
			index.drawOutput(new Uint8ClampedArray(responsePayload.data), responsePayload.width, responsePayload.height)
		})
	},
	cancelNewImg: function() {
		currentRequestVersion++;
		const request = { name: 'cancel-new', payload: { 
			Version: currentRequestVersion,
		}};
		astilectron.sendMessage(request, function(response) {
			if (response.payload.version === currentRequestVersion) {
				index.finishLoading();
			} 
			if (response.name === 'error') {
				asticode.notifier.error(response.payload);
                return
			}
		})
	},
	initPalette: function() {
		const colorCtr = document.getElementById('color-sample-ctr')
		Object.keys(colorMap)
			.forEach(color => {
				const el = document.createElement('div')
				el.id = color
				el.classList.add('color-sample')
				el.classList.add('shadowed')
				el.style.background = colorMap[color]
				colorCtr.appendChild(el)
				el.addEventListener('click', e => {
					activeColor = el.id;
					document.querySelectorAll('.color-sample')
						.forEach(cc => {
							cc.classList.remove('active')
						})
					el.classList.add('active')
				})
			})

		const black = document.getElementById('black')
		black.classList.add('active')
		activeColor = 'black'
	},
	initInput: function() {
		const canvas = document.getElementById('input')
		canvas.width = inputWidth
		canvas.height = inputHeight
		const ctx = canvas.getContext('2d')
		ctx.fillStyle = colorMap['white'];
		ctx.fillRect(0, 0, canvas.width, canvas.height)
		let isMouseDown = false

		function paintPixel(e) {
			if (isMouseDown) {
				const rect = canvas.getBoundingClientRect()
				const loc = {
					x: e.clientX - rect.left,
					y: e.clientY - rect.top
				}
				
				ctx.fillStyle = colorMap[activeColor]
				ctx.fillRect(loc.x - (loc.x % inputPixelSize), loc.y - (loc.y % inputPixelSize), inputPixelSize, inputPixelSize)
			}
		}

		canvas.addEventListener('mousemove', function(e) {
			paintPixel(e)
		})

		canvas.addEventListener('mousedown', function(e) {
			isMouseDown = true
			paintPixel(e)
		})

		document.body.addEventListener('mouseup', function(e) {
			if (isMouseDown) {
				isMouseDown = false
				index.requestNewImg();
			}
		})
	},
	clearInput: function() {
		const canvas = document.getElementById('input')
		canvas.width = inputWidth
		canvas.height = inputHeight
		const ctx = canvas.getContext('2d')
		ctx.fillStyle = colorMap['white'];
		ctx.fillRect(0, 0, canvas.width, canvas.height)
		index.requestNewImg();
	},
	getInputData: function() {
		const canvas = document.getElementById('input')
		const ctx = canvas.getContext('2d')
		const imgData = ctx.getImageData(0, 0, canvas.width, canvas.height)
		return new ImageData(imgData.data.filter((d, i) => {
			return i % (4 * inputPixelSize) < 4 && i % (4 * inputPixelSize * inputWidth) < 4 * inputWidth
		}), inputWidth / inputPixelSize, inputHeight / inputPixelSize)
	},
	loadInput: function(path, cb) {
		const img = document.createElement('img');
		img.src = path;
		img.onload = function(e) {
			// Set image in input canvas
			const canvas = document.getElementById('input')
			canvas.width = this.width
			canvas.height = this.height
			const ctx = canvas.getContext('2d')
			ctx.drawImage(this, 0, 0)
			// Extract image data
			if (cb) {
				cb(ctx.getImageData(0, 0, this.width, this.height))
			}
		}
	},
	drawOutputFromFile: function(path) {
		const img = document.createElement('img');
		img.src = path;
		img.onload = function(e, cb) {
			// Set image in output canvas
			const canvas = document.getElementById('output')
			canvas.width = this.width
			canvas.height = this.height
			const ctx = canvas.getContext('2d')
			ctx.drawImage(this, 0, 0)
			// Extract image data
			if (cb) {
				cb(ctx.getImageData(0, 0, this.width, this.height))
			}
		}
	},
	drawOutput: function(data, width, height, cb) {
		let imageData = new ImageData(data, width, height)
		imageData = scaleImageData(imageData, outputPixelSize)
		const canvas = document.getElementById('output')
		canvas.width = width * outputPixelSize * 2
		canvas.height = height * outputPixelSize * 2
		const ctx = canvas.getContext('2d')
		ctx.putImageData(imageData, 0, 0)
		ctx.putImageData(imageData, width * outputPixelSize, 0)
		ctx.putImageData(imageData, 0, height * outputPixelSize)
		ctx.putImageData(imageData, width * outputPixelSize, height * outputPixelSize)

		// Extract image data
		if (cb) {
			cb(ctx.getImageData(0, 0, width, height))
		}
	},
}

function scaleImageData(imageData, scale) {
	const newData = []
	for (let i = 0; i < imageData.data.length; i += 4) {
		for (let n = 0; n < scale; n++) {
			newData.push(...[imageData.data[i], imageData.data[i+1], imageData.data[i+2], imageData.data[i+3]])
		}
		if ((i / 4) % imageData.width === imageData.width - 1) {
			const rowCopy = newData.slice(-imageData.width * scale * 4)
			for (let n = 0; n < scale - 1; n++) {
				newData.push(...rowCopy)
			}
		}
	}
	return new ImageData(new Uint8ClampedArray(newData), imageData.width * scale, imageData.height * scale)
}


