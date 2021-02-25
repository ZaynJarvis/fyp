/*
 * MinIO Cloud Storage (C) 2020 MinIO, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from "react"
import { Modal, ModalHeader, ModalBody } from "react-bootstrap"
import { a11yLight, CopyBlock } from "react-code-blocks";
import axios from 'axios';

class PreviewObjectModal extends React.Component {
  constructor(props) {
    super(props)
    this.state = {
      url: "",
      res: null,
    }
  }

  componentDidMount() {
    this.props.getObjectInfo(this.props.object.name, (res) => {
      this.setState({
        ...this.state,
        res: res
      })
    })
    this.props.getObjectURL(this.props.object.name, (url) => {
      this.setState({
        ...this.state,
        url: url,
      })
    })
  }

  render() {
    const { hidePreviewModal, object } = this.props
    return (
      <Modal
        show={true}
        animation={false}
        onHide={hidePreviewModal}
        bsSize="large"
      >
        <ModalHeader>{object.name}</ModalHeader>
        <ModalBody>
          <div className="input-group">
            {this.state.url && (
              <object data={this.state.url} style={{ display: "block", width: "100%" }}>
                <h3 style={{ textAlign: "center", display: "block", width: "100%" }}>
                  Do not have read permissions to preview "{this.props.object.name}"
                </h3>
              </object>
            )}
          </div>
          <div>
            {this.state.res && (<div
              style={{
                width: '100%',
                flex: 1,
                paddingBottom: '1em',
              }}
            >
              <h2 style={{ textAlign: 'center' }}>Result</h2>
              <CopyBlock
                text={JSON.stringify(this.state.res.result, "", "  ")}
                theme={a11yLight}
                wrapLines={false}
              />
              {
                this.state.res.text && (
                  <div>
                    <h2 style={{ textAlign: 'center' }}>Collection Notes</h2>
                    <p>{this.state.res.text}</p>
                  </div>
                )
              }
            </div>)}
          </div>
        </ModalBody>
      </Modal>
    )
  }
}
export default PreviewObjectModal
