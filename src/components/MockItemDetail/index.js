import React from 'react';
import InputGroup from 'react-bootstrap/InputGroup';
import Col from 'react-bootstrap/Col';
import BigTextInput from '../BigTextInput'
import MockItem from '../MockItem'

const MockItemDetail = (props) => {
    if(!props.data.method) {
        return (
            <h3 className="header">Choose a template from left</h3>
        )
    }
    let method = props.data.method.toString().toUpperCase();
    return (
        <Col style={{ alignSelf: 'center' }}>
            <MockItem data={props.data} ></MockItem>
            {method === 'POST' ? 
            <BigTextInput label={'Request'} disabled value={JSON.stringify(props.data.request)} ></BigTextInput>
            :
            null
            }
            <BigTextInput label={'Reponse'}  disabled value={JSON.stringify(props.data.response)}></BigTextInput>
        </Col>
    );
};

export default MockItemDetail;
