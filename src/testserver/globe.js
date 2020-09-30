var viewer = new Cesium.Viewer("cesiumContainer", {
    selectionIndicator: false,
    infoBox: false,
    terrainProvider: Cesium.createWorldTerrain(),
});

viewer.camera.lookAt(
    Cesium.Cartesian3.fromDegrees(0, 0, 0.0),
    new Cesium.Cartesian3(0.0, 0.0, 42000000.0)
);

viewer.cesiumWidget.screenSpaceEventHandler.removeInputAction(
    Cesium.ScreenSpaceEventType.LEFT_DOUBLE_CLICK
);

var handler = new Cesium.ScreenSpaceEventHandler(viewer.canvas);
handler.setInputAction(onLeftMouseClick, Cesium.ScreenSpaceEventType.LEFT_CLICK);
viewer.camera.lookAtTransform(Cesium.Matrix4.IDENTITY);
viewer.timeline.destroy()
viewer.sceneModePicker.destroy()
viewer.navigationHelpButton.destroy()
viewer.homeButton.destroy()
viewer.geocoder.destroy()
viewer.animation.destroy()

var data = {
    startLat: "",
    startLng: ""
}

var routes = [
    [],
    [],
    [],
    [],
    []
]

if (!viewer.scene.pickPositionSupported) {
    window.alert("This browser does not support pickPosition.");
}

function routeTest(jsonData) {
    if (jsonData != "false") {
        tempData = JSON.parse(jsonData)
        for (j = 0; j < tempData.length; j++) {
            for (i = 0; i < tempData[j].features.length; i++) {
                var coord1d = []
                var coordinates = tempData[j].features[i].geometry.coordinates;
                coordinates.forEach(element => coord1d.push(element[0], element[1]))
                var color = Cesium.Color.RED
                var width = 4
                switch (j) {
                    case 0:
                        color = Cesium.Color.CRIMSON
                        break
                    case 1:
                        color = Cesium.Color.GREENYELLOW
                        break
                    case 2:
                        color = Cesium.Color.ORANGE
                        break
                    case 3:
                        color = Cesium.Color.PURPLE
                        break
                    case 4:
                        color = Cesium.Color.ROYALBLUE
                        width = 2
                        break
                }
                routes[j][i] = viewer.entities.add({
                    name: j,
                    positions: coord1d.slice(0, 2),
                    polyline: {
                        positions: Cesium.Cartesian3.fromDegreesArray(coord1d),
                        width: width,
                        material: color,
                    }
                })
            }
        }

    }
}

function toggle(cb, i) {
    if (cb.checked) {
        console.log(routes[i])
        for (var j = 0; j < routes[i].length; j++) {
            routes[i][j].show = true
        }
    } else {
        for (var j = 0; j < routes[i].length; j++) {
            routes[i][j].show = false
        }
    }
}

function createPoint(worldPosition, start = false) {
    var point

    var text = "End"
    console.log(worldPosition)
    point = viewer.entities.add({
        position: worldPosition,
        point: {
            color: Cesium.Color.WHITE,
            pixelSize: 5,
        },
    });

    if (start == true) {
        text = "Start"
    }

    point = viewer.entities.add({
        position: worldPosition,
        label: {
            height: 20000000,
            text: text,
            font: '14pt monospace',
            style: Cesium.LabelStyle.FILL_AND_OUTLINE,
            outlineWidth: 2,
            verticalOrigin: Cesium.VerticalOrigin.TOP,
            pixelOffset: new Cesium.Cartesian2(0, 32),
            eyeOffset: new Cesium.Cartesian3(0, 0, -3000000)
        }
    });
    return point;
}

function nbTest(testData) {
    var td = JSON.parse(testData)
    console.log(td.Point)
    createColoredPoint(Cesium.Cartesian3.fromDegrees(td.Point[0], td.Point[1]), Cesium.Color.CHARTREUSE)
    if (document.getElementById("simpleNbs").checked) {
        if (td.Nnbs != null) {
            for (i = 0; i < td.Nnbs.length; i++) {
                createColoredPoint(Cesium.Cartesian3.fromDegrees(td.Nnbs[i][0], td.Nnbs[i][1]), Cesium.Color.ALICEBLUE)
            }
        }
    } else {
        for (i = 0; i < td.Nbs.length; i++) {
            createColoredPoint(Cesium.Cartesian3.fromDegrees(td.Nbs[i][0], td.Nbs[i][1]), Cesium.Color.RED)
        }
    }
}

function onLeftMouseClick(event) {
    // We use `viewer.scene.pickPosition` here instead of `viewer.camera.pickEllipsoid` so that
    // we get the correct point when mousing over terrain.
    var earthPosition = viewer.scene.pickPosition(event.position);

    // `earthPosition` will be undefined if our mouse is not over the globe.
    if (Cesium.defined(earthPosition)) {

        const cartographic = viewer.scene.globe.ellipsoid.cartesianToCartographic(earthPosition);
        const longitudeString = Cesium.Math.toDegrees(cartographic.longitude).toFixed(15);
        const latitudeString = Cesium.Math.toDegrees(cartographic.latitude).toFixed(15);

        data["startLat"] = latitudeString
        data["startLng"] = longitudeString
        $.ajax({
            url: "/point",
            data: data,
            earthPosition: { ep: [longitudeString, latitudeString] }
        }).done(nbTest);
    }
}

function createColoredPoint(worldPosition, color) {
    var point = viewer.entities.add({
        position: worldPosition,
        point: {
            color: color,
            pixelSize: 6,
        },
    });
    return point;
}

$("#clearButton").click(function() {
    viewer.entities.removeAll();
});

$("#gridButton").click(function() {
    $.ajax({
        url: "/grid",
        data: data
    }).done(function(data) {
        if (data != "false") {
            points = JSON.parse(data)
            for (i = 0; i < points.length; i++) {
                createColoredPoint(Cesium.Cartesian3.fromDegrees(points[i][0], points[i][1]), Cesium.Color.ALICEBLUE)
            }
        }
    });
});

$("#showNbs").click(function() {
    $.ajax({
        url: "/id",
        data: { id: document.getElementById("gridIDX").value }
    }).done(function(data) {
        if (data != "false") {
            console.log(data)
            nbTest(data)
        }
    });
});

$("#showRoutes").click(function() {
    x = document.getElementById("gridIDXx").value
    y = document.getElementById("gridIDXy").value
    data = {
        gridIDXx: document.getElementById("gridIDXx").value,
        gridIDXy: document.getElementById("gridIDXy").value
    }

    $.ajax({
        url: "/startend",
        data: data
    }).done(function(data) {
        data = JSON.parse(data)
        if (data != "false") {
            console.log(data)
            createPoint(Cesium.Cartesian3.fromDegrees(data[0][0], data[0][1]), true);
            createPoint(Cesium.Cartesian3.fromDegrees(data[1][0], data[1][1]));
        }
    });

    $.ajax({
        url: "/route",
        data: data,
        earthPosition: { ep: [x, y] }
    }).done(routeTest);
});