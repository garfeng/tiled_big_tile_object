#  Generate Big Objects for Tiled

[中文版](translations/zh_cn.md) | [Download](https://github.com/garfeng/tiled_big_tile_object/releases)

There are some item with size more than 1 x 1. It is difficult to place them in map as `objects`.

![image-20221103214954877](README.assets/image-20221103214954877.png)



Put them to a map (tile layer) separate from each other.

![image-20221103220438444](README.assets/image-20221103220438444.png)



Export the map as Image.

![image-20221103220518063](README.assets/image-20221103220518063.png)



Run this tool. It will group objects with same size to one image.

| objects_1x1.png                      | objects_1x2.png                      | objects_2x2.png                      |
| ------------------------------------ | ------------------------------------ | ------------------------------------ |
| ![1x1](examples/dst/objects_1x1.png) | ![1x2](examples/dst/objects_1x2.png) | ![2x2](examples/dst/objects_2x2.png) |



Then you could create new Tilesets in Tiled with their tile size.

![image-20221103221258527](README.assets/image-20221103221258527.png)



Draw your map with then.  [Example maps](./examples).

![image-20221103225551089](README.assets/image-20221103225551089.png)





