import React from 'react';
import InputGroup from 'react-bootstrap/InputGroup';
import FormControl from 'react-bootstrap/FormControl';

const PrefixedInput = (props) => {
    return (
        <InputGroup className="mb-3">
            <InputGroup.Prepend>
                <InputGroup.Text id="basic-addon3">
                https://localhost:3000/
                </InputGroup.Text>
            </InputGroup.Prepend>
            <FormControl
                id="basic-url"
                aria-describedby="basic-addon3"
                value={props.value}
                onChange={props.onChange}

            />
        </InputGroup>

    );
};

export default PrefixedInput;

{/* <input type="text" value={props.value} onChange={props.onChange} /> */}

// formlar düzenleencek. validasonlar formlarla yapılcak