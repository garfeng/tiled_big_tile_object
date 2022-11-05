import {Component, useState} from 'react';
import './App.css';
import {Greet} from "../wailsjs/go/main/App";

import {Col, Layout, Row, Input,Form, Button, InputNumber, Typography, Divider} from "antd"
import { Footer } from 'antd/lib/layout/layout';

import {GithubOutlined} from '@ant-design/icons';

const {Content} = Layout;
const {TextArea} = Input;
const {Title, Link} = Typography;

class App extends Component {
    constructor(props) {
        super(props);
        this.state = {
            inputImagesStr:"",
            inputImages : [],
            tileSize : 48,
            dstWidth : 640,
            dstHeight: 640,
            dstRoot : "./dst",
            dstPrefix : "objects",

        }
    }
    onFinish = (e) => {

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
                            <Form name="basic" labelCol={{span:4}} wrapperCol={{span:16}} onFinish={this.onFinish}>
                                <Form.Item label="Input images" name="inputImage">
                                    <Input.Group>
                                    <Input readOnly={true} value={this.state.inputImagesStr} style={{ width: 'calc(100% - 100px)' }}/>
                                    <Button type='primary'>Select</Button>
                                    </Input.Group>
                                </Form.Item>
                                <Form.Item label="Tile size" name="tileSize">
                                    <InputNumber defaultValue={this.state.tileSize} value={this.state.tileSize} addonAfter={"px"}></InputNumber>
                                </Form.Item>
                                <Form.Item label="Dst Size" name="dstSize">
                                    <Input.Group>
                                    <InputNumber defaultValue={this.state.dstWidth} value={this.state.dstWidth} style={{width:"5rem"}}></InputNumber> {" x "}
                                    <InputNumber defaultValue={this.state.dstWidth} value={this.state.dstHeight} style={{width:"7rem"}} addonAfter={"px"}></InputNumber>
                                    </Input.Group>
                                </Form.Item>
                                <Form.Item label="Dst root" name="dstRoot">
                                    <Input.Group compact>
                                    <Input value={this.state.dstRoot} readOnly={true} style={{ width: 'calc(100% - 100px)' }}></Input>
                                    <Button type='primary'>Select</Button>
                                    </Input.Group>
                                </Form.Item>
                                <Form.Item label="Dst prefix" name="dstPrefix">
                                    <Input defaultValue={this.state.dstPrefix} value={this.state.dstPrefix} style={{ width: 'calc(100% - 100px)' }}></Input>
                                </Form.Item>
                                <Form.Item label="Generate">
                                    <Button type='primary'>Generate</Button>
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
