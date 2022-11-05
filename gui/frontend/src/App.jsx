import {Component, useState} from 'react';
import './App.css';
import {SelectImages, SelectDstRoot, Generate} from "../wailsjs/go/main/App";
import {Col, Layout, Row, Input,Form, Button, InputNumber, Typography, Divider, message} from "antd"
import { Footer } from 'antd/lib/layout/layout';

import {GithubOutlined} from '@ant-design/icons';

const {Content} = Layout;
const {TextArea} = Input;
const {Title, Link} = Typography;

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            inputImages : [],
            tileSize : 48,
            dstWidth : 640,
            dstHeight: 640,
            dstRoot : "./dst",
            dstPrefix : "objects",
        }
    }

    SelectImages = (e) => {
        SelectImages().then((images) => {
            if(images.length > 0) {
                this.setState({
                    inputImages : images,
                })
            }
        }).catch(
            (reason) => {
                // do nothing
            }
        )
    }

    SelectDstRoot = (e) => {
        SelectDstRoot().then((dstRoot)=>{
            if (dstRoot.length > 0) {
                this.setState(
                    {
                        dstRoot: dstRoot,
                    }
                )
            }

        }).catch((reason) => {
            // do nothing
        });
    }

    Generate = (e) => {
        var param = {
            tileSize: this.state.tileSize,
            dstWidth: this.state.dstWidth,
            dstHeight: this.state.dstHeight,
            srcImages: this.state.inputImages,
            dstRoot: this.state.dstRoot,
            dstPrefix: this.state.dstPrefix,
        };
        console.log(param)
        Generate(param).then( (err)=> {
            if (err != null) {
                message.error(err.message);
            } else {
                message.success("finished");
            }
        }  ).catch( reason => {
            message.error(reason);
        } )
    }

    setValue = (name) => {
        return (value) => {
            if (typeof value == "string" || typeof value == "number") {
                this.setState(
                    {
                        [name]: value
                    }
                )
            } else {
                this.setState(
                    {
                        [name]: value.target.value
                    }
                )
            }
        }
    }

    render() {

        return (
            <Layout style={{width:"100%", height:"100%"}}>
            <Content>
                <Row>
                    <Col span={20} offset={2}>
                                <div style={{height:"2em"}}></div>
                                <Col offset={4}>
                                <Title>
                                Generate Tiled Big Objects
                            </Title>
                                </Col>
                                <Divider/>
                            <Form name="basic" labelCol={{span:4}} wrapperCol={{span:16}}>

                                <Form.Item label="Input images" name="inputImages">
                                    <Input.Group>
                                    <Input readOnly={true} value={this.state.inputImages.join(";")} style={{ width: 'calc(100% - 100px)' }}/>
                                    <Button type='primary' onClick={this.SelectImages}>Select</Button>
                                    </Input.Group>
                                </Form.Item>

                                <Form.Item label="Tile size" name="tileSize" initialValue={this.state.tileSize}>
                                    <InputNumber onChange={this.setValue("tileSize")}  value={this.state.tileSize} addonAfter={"px"}></InputNumber >
                                </Form.Item>

                                <Form.Item label="Dst Size" name="dstSize">
                                    <Input.Group>
                                    <InputNumber onChange={this.setValue("dstWidth")} value={this.state.dstWidth} style={{width:"5rem"}}></InputNumber> {" x "}
                                    <InputNumber onChange={this.setValue("dstHeight")} value={this.state.dstHeight} style={{width:"7rem"}} addonAfter={"px"}></InputNumber>
                                    <span>{" "}Will auto set to <code>TileSize x N</code></span>
                                    </Input.Group>
                                </Form.Item>

                                <Form.Item label="Dst root" name="dstRoot">
                                    <Input.Group compact>
                                    <Input value={this.state.dstRoot} readOnly={true} style={{ width: 'calc(100% - 100px)' }}></Input>
                                    <Button type='primary' onClick={this.SelectDstRoot}>Select</Button>
                                    </Input.Group>
                                </Form.Item>

                                <Form.Item label="Dst prefix" name="dstPrefix" initialValue={this.state.dstPrefix}>
                                    <Input onChange={this.setValue("dstPrefix")} value={this.state.dstPrefix} style={{ width: 'calc(100% - 100px)' }}></Input>
                                </Form.Item>

                                <Form.Item label="Generate">
                                    <Button type='primary' onClick={this.Generate}>Generate</Button>
                                </Form.Item>
                            </Form>
                    </Col>
                </Row>
            </Content>
            <Footer style={{textAlign:"center"}}> <Link href='https://github.com/garfeng/tiled_big_tile_object'><GithubOutlined /> Github</Link> |
                Driven by <Link href='https://github.com/wailsapp/wails'>Wails</Link> (Create beautiful applications using Go) | <Link href='https://rpg.blue'>Project 1</Link>
            </Footer>
          </Layout>
        )
    }
}

export default App
