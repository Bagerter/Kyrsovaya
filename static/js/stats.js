document.addEventListener("DOMContentLoaded", function() {
    // Карта посещений
    am4core.ready(function() {
        // Создание карты
        let map = am4core.create("visitsMap", am4maps.MapChart);
        map.geodata = am4geodata_worldLow;
        map.projection = new am4maps.projections.Miller();

        let polygonSeries = map.series.push(new am4maps.MapPolygonSeries());
        polygonSeries.useGeodata = true;
        
        // Данные для демонстрации
        polygonSeries.data = [
            { id: "RU", value: 400, fill: am4core.color("#F05") }, // Пример данных
            // Добавьте больше стран и значений
        ];

        let polygonTemplate = polygonSeries.mapPolygons.template;
        polygonTemplate.tooltipText = "{name}: {value}";
        polygonTemplate.fill = am4core.color("#74B266");

        polygonSeries.heatRules.push({
            property: "fill",
            target: polygonSeries.mapPolygons.template,
            min: am4core.color("#FFD1AA"),
            max: am4core.color("#7B3625")
        });

        map.series.push(polygonSeries);
		map.zoomControl = new am4maps.ZoomControl();
    });

    am4core.ready(function() {
    // Создание карты атак
    let attacksMap = am4core.create("attacksMap", am4maps.MapChart);

    // Загрузка геоданных для карты мира
    attacksMap.geodata = am4geodata_worldLow;
    attacksMap.projection = new am4maps.projections.Miller();

    let polygonSeriesAttacks = attacksMap.series.push(new am4maps.MapPolygonSeries());
    polygonSeriesAttacks.useGeodata = true;

    // Демонстрационные данные для атак по странам
    polygonSeriesAttacks.data = [
        { id: "US", value: 200, fill: am4core.color("#F00") }, // Пример данных
        // Добавьте больше данных о странах и атаках
    ];

    let polygonTemplateAttacks = polygonSeriesAttacks.mapPolygons.template;
    polygonTemplateAttacks.tooltipText = "{name}: {value} атак";
    polygonTemplateAttacks.fill = am4core.color("#74B266");

    // Правила для изменения цвета в зависимости от количества атак
    polygonSeriesAttacks.heatRules.push({
        property: "fill",
        target: polygonSeriesAttacks.mapPolygons.template,
        min: am4core.color("#FFD1AA"),
        max: am4core.color("#7B3625")
    });

    // Добавление zoomControl для удобства навигации по карте
    attacksMap.zoomControl = new am4maps.ZoomControl();
});

});
