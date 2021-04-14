import React from "react"
import { Modal } from "react-bootstrap"
import logo from "../../img/logo.svg"
import { FormGroup, FormControl, ControlLabel, Button } from 'react-bootstrap'
import { useFormik } from 'formik';

export const ConfigModal = ({ hideConfig }) => {
  const [svc, setSvc] = React.useState('no');
  const [record, setRecord] = React.useState([]);

  React.useEffect(() => {
    fetch(`http://localhost:8900/services`, { method: 'GET' })
      .then(
        function (response) {
          if (response.status !== 200) {
            console.log('Looks like there was a problem. Status Code: ' +
              response.status);
            return;
          }

          response.json().then(data => {
            setRecord(data)
          });
        }
      )
      .catch(function (err) {
        console.log('Fetch Error :-S', err);
      });
  }, [])
  return (
    <Modal
      className="modal-config modal-dark"
      animation={false}
      show={true}
      onHide={hideConfig}
    >
      <button className="close" onClick={hideConfig}>
        <span>×</span>
      </button>
      <div className="mc-inner">
        <div className="mc-list">
          {record.map((v, i) => <Card key={i} svc={v} me={svc === v} onClick={() => setSvc(v)} />)}
        </div>
        <div className="mc-form">
          {svc === 'no' ? <div className="mc-cover"><img className="maii-logo" src={logo} alt="" /></div> :
            <RuleForm svc={svc} />}
        </div>
      </div>
    </Modal>
  )
}

const Card = ({ svc, me, onClick }) => {
  return (
    <div className='mc-card' onClick={onClick}>
      {svc} <div style={{ backgroundColor: me ? 'lightgreen' : 'transparent' }} />
    </div>
  )
}

const opMap = {
  1: 'exist',
  2: 'not_exist',
  5: 'lt',
  7: 'gt',
}

const RuleForm = ({ svc }) => {
  const obj = {
    version: "1",
    service: svc,
    object_storage_path: "minio:9000",
    document_storage_path: "mongo:27017",
  }
  const [status, setStatus] = React.useState(0)

  React.useEffect(() => {
    setStatus(0)
    fetch(`http://localhost:8900/services/${svc}`, { method: 'GET' })
      .then(
        function (response) {
          if (response.status !== 200) {
            console.log('Looks like there was a problem. Status Code: ' +
              response.status);
            return;
          }

          response.json().then(data => {
            if (data.rules && data.rules.length > 0) {
              rule.setFieldValue('field', data.rules[0].field)
              rule.setFieldValue('op', opMap[data.rules[0].op])
              rule.setFieldValue('operand', data.rules[0].operand)
              rule.setFieldValue('sample_rate', data.rules[0].sample_rate)
            } else {
              rule.setFieldValue('field', '')
              rule.setFieldValue('op', '')
              rule.setFieldValue('operand', '')
              rule.setFieldValue('sample_rate', 1)
            }
          });
        }
      )
      .catch(function (err) {
        console.log('Fetch Error :-S', err);
      });
  }, [svc])

  let request = (body) => {
    console.log('body: ' + JSON.stringify(body))
    fetch(`http://localhost:8900/services/${body.service}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })
      .then(response => response.json())
      .then(() => {
        setStatus(1)
      })
      .catch((error) => {
        setStatus(0)
        alert('Error: ' + error)
      });
  }

  const rule = useFormik({
    initialValues: {
      field: '',
      op: '',
      operand: '',
      sample_rate: 1,
    },
    onReset: rule => {
      rule.operand = ''
      rule.op = ''
      rule.field = ''
      request({ ...obj, rules: [] })
    },
    onSubmit: rule => {
      if (rule.field === '') {
        alert('Error: field should not be null')
        return
      } else if (rule.op === '') {
        alert('Error: operator should not be null')
        return
      } else if (rule.op === 'exist' || rule.op === 'not_exist') {
        rule.operand = ''
      } else if (rule.operand === '') {
        alert('Error: operand should not be null')
        return
      }
      request({ ...obj, rules: [rule] })
    }
  })

  return (<form onSubmit={rule.handleSubmit} onReset={rule.handleReset}>
    <FormGroup controlId='field'>
      <ControlLabel>Form Field</ControlLabel>
      <FormControl help="field..."
        value={rule.values.field}
        onChange={rule.handleChange} />
    </FormGroup>
    <FormGroup controlId="op">
      <ControlLabel>Operator</ControlLabel>
      <FormControl componentClass="select" placeholder="select"
        value={rule.values.op}
        onChange={rule.handleChange}
      >
        <option value='' disabled>select operator</option>
        <option value='exist'>Exists</option>
        <option value='not_exist'>Not Exists</option>
        <option value='lt'>Less Than</option>
        <option value='gt'>Greater Than</option>
      </FormControl>
    </FormGroup>
    {
      ['lt', 'gt'].includes(rule.values.op) &&
      <FormGroup controlId="operand">
        <ControlLabel>Operand</ControlLabel>
        <FormControl type="text"
          value={rule.values.operand}
          onChange={rule.handleChange}
        />
      </FormGroup>
    }
    <Button type="reset">Reset</Button>{' '}
    <Button type="submit">Configure</Button>
    <span>{status === 1 ? 'Success!' : ''}</span>
  </form>);
}

export default ConfigModal

