import UIKit

struct UserInfoResponse: Decodable {
    let error_text: String
    let has_error: Bool
    let message: UserInfo
}

struct AllTagsResponse: Decodable {
    let error_text: String
    let has_error: Bool
    let message: [AllTags]
}

struct AllTags: Decodable, Encodable {
    let id: Int64
    let title: String
}

struct UserInfo: Decodable {
    let login: String
    let email: String
    let name: String
    let age: String
    let f_tags: String
    let uf_tags: String
}

class ModelDataTags {
    
    var allTags = [AllTags]()
    var unfavoriteTags = [AllTags]()
    
    init() {
        
    }
    
    func loadData(all:[AllTags], my:[AllTags]) {
        setupData(all: all, my: my)
    }
    
    func setupData(all: [AllTags], my: [AllTags]) {
        allTags = all
        unfavoriteTags = my
    }
    
    func deleteSameTags() {
        var i = 0
        for alltag in allTags {
            for untag in unfavoriteTags {
                if alltag.id == untag.id {
                    allTags.remove(at: i)
                    i -= 1
                }
            }
            i += 1
        }
    }
    
    func addElement(index: Int) {
        if allTags.count != 0 {
            unfavoriteTags.append(allTags[index])
            allTags.remove(at: index)
        }
    }
    
    func deleteElement(index: Int) {
        if unfavoriteTags.count != 0 {
            allTags.append(unfavoriteTags[index])
            unfavoriteTags.remove(at: index)
        }
        
    }
}

class MainBarViewController: UITabBarController {
    
    
    @IBOutlet weak var qwerty: UITabBar!
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }

}

class AllEventsViewCOntroller: UIViewController {
    
    override func viewDidAppear(_ animated: Bool) {
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
            
        }
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
}

class EventsViewCOntroller: UIViewController {
    
    override func viewDidAppear(_ animated: Bool) {
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
        }
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }
}


class ProfileTagsViewController: UIViewController {
    
    @IBOutlet weak var pickerView: UIPickerView!
    @IBOutlet weak var AddTagButton: UIButton!
    @IBOutlet weak var DeleteTagButton: UIButton!
    @IBOutlet weak var UnfavoriteTagsCollection: UIPickerView!
    
    private var ResAllTags = [AllTags]()
    private var ResUnfavTags = [AllTags]()
    var modelData = ModelDataTags()
    
    override func viewDidAppear(_ animated: Bool) {
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
        }
        
    }
    
    func updateViewDisplay() {
        modelData.deleteSameTags()
        ResAllTags = modelData.allTags
        ResUnfavTags = modelData.unfavoriteTags
        print(ResAllTags)
        print(ResUnfavTags)
        pickerView.dataSource = self
        pickerView.delegate = self
        UnfavoriteTagsCollection.dataSource = self
        UnfavoriteTagsCollection.delegate = self
    }
    
    @IBAction func DeleteTagButton(_ sender: UIButton) {
        let selected = UnfavoriteTagsCollection.selectedRow(inComponent: 0)
        modelData.deleteElement(index: selected)
        updateViewDisplay()
    }
    
    @IBAction func AddTagButton(_ sender: UIButton) {
        let selected = UnfavoriteTagsCollection.selectedRow(inComponent: 0)
        modelData.addElement(index: selected)
        updateViewDisplay()
    }
    
    func getUnfavoriteTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/get_unfavorite_tags") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(AllTagsResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        DispatchQueue.main.async {
                            ResUnfavTags = httpResp.message
                            modelData.loadData(all: ResAllTags, my: ResUnfavTags)
                            updateViewDisplay()
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    func updateUnfavoriteTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/update_unfavorite_tags") else { return }
        var request = URLRequest(url: url)
        let encoder = JSONEncoder()
        let param = try! encoder.encode(ResUnfavTags)
        print(String(data: param, encoding: .utf8)!)

        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        //guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = param
        
        let session = URLSession.shared
        session.dataTask(with: request) { [] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(HTTPResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                    }
                    
                } else {
                    print(httpResp.error_text)
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    func getAllTags() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/tags/") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(AllTagsResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        DispatchQueue.main.async {
                            ResAllTags = httpResp.message
                            modelData.loadData(all: ResAllTags, my: ResUnfavTags)
                            updateViewDisplay()
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                    let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                    present(newVC, animated: true, completion: nil)
                }
            } catch {
                print(error)
            }
        }.resume()
    }

    
    override func viewDidDisappear(_ animated: Bool) {
        updateUnfavoriteTags()
        navigationController?.popViewController(animated: true)
        dismiss(animated: true, completion: nil)
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
        pickerView.tag = 1
        UnfavoriteTagsCollection.tag = 2
        getAllTags()
        getUnfavoriteTags()
    }

}

extension ProfileTagsViewController: UIPickerViewDataSource {
    func pickerView(_ pickerView: UIPickerView, numberOfRowsInComponent component: Int) -> Int {
        if pickerView.tag == 1 {
            return modelData.allTags.count
        } else if pickerView.tag == 2 {
            return modelData.unfavoriteTags.count
        } else {
            return 0
        }
    }
    
    func numberOfComponents(in pickerView: UIPickerView) -> Int {
        return 1
    }
}

extension ProfileTagsViewController: UIPickerViewDelegate {
    func pickerView(_ pickerView: UIPickerView, titleForRow row: Int, forComponent component: Int) -> String? {
        if pickerView.tag == 1 {
            let a = modelData.allTags[row]
            return a.title
        } else if pickerView.tag == 2 {
            let a = modelData.unfavoriteTags[row]
            return a.title
        } else {
            return ""
        }
    }
}

class ProfileViewController: UIViewController {

    @IBOutlet weak var profileName: UILabel!
    @IBOutlet weak var profileLogin: UILabel!
    @IBOutlet weak var profileEmail: UILabel!
    @IBOutlet weak var profileAge: UILabel!
    @IBOutlet weak var favoriteTags: UITextView!
    @IBOutlet weak var unfavoriteTags: UITextView!
    @IBOutlet weak var oldPassword: UITextField!
    @IBOutlet weak var newPassword: UITextField!
    @IBOutlet weak var changePasButton: UIButton!
    @IBOutlet weak var errorLabel: UILabel!
    @IBOutlet weak var exitButton: UIButton!
    
    override func viewDidAppear(_ animated: Bool) {
        DispatchQueue.main.async {
            self.newPassword.text = ""
            self.oldPassword.text = ""
            self.errorLabel.text = ""
        }
        let Token = UserDefaults.standard
        if Token.string(forKey: "Token") == "" {
            tabBarController?.selectedIndex = 0
            let storyBoard = UIStoryboard(name: "Main", bundle: nil)
            let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
            present(newVC, animated: true, completion: nil)
        } else {
            updateUserInfo()
        }
    }
    
    func updateUserInfo() {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/auth/me") else { return }
        var request = URLRequest(url: url)
        let param = ["q": "q"]
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(UserInfoResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        DispatchQueue.main.async {
                            profileName.text = httpResp.message.name
                            profileLogin.text = httpResp.message.login
                            profileEmail.text = httpResp.message.email
                            profileAge.text = httpResp.message.age
                            favoriteTags.text = httpResp.message.f_tags
                            unfavoriteTags.text = httpResp.message.uf_tags
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        tabBarController?.selectedIndex = 0
                        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
                        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
                        present(newVC, animated: true, completion: nil)
                    }
                }
            } catch {
                print(error)
            }
        }.resume()
        
    }
    
    override func viewDidLoad() {
        super.viewDidLoad()
    }

    @IBAction func changePasButton(_ sender: UIButton) {
        let Token = UserDefaults.standard
        guard let url = URL(string: "http://127.0.0.1:8088/api/auth/change_pass") else { return }
        var request = URLRequest(url: url)
        let param = ["old_password": oldPassword.text!, "new_password": newPassword.text!]
        print(param)
        request.httpMethod = "POST"
        request.addValue("application/json", forHTTPHeaderField: "Content-Type")
        request.addValue(Token.string(forKey: "Token") ?? "", forHTTPHeaderField: "Token")
        guard let httpBody = try? JSONSerialization.data(withJSONObject: param, options: []) else { return }
        request.httpBody = httpBody
        
        let session = URLSession.shared
        session.dataTask(with: request) { [self] (data, response, error) in
            if let response = response {
                print(response)
            }
            guard let data = data else {
                return
            }
            
            do {
                let httpResp = try JSONDecoder().decode(HTTPResponse.self, from: data)
                print(httpResp.message)
                if httpResp.has_error == false {
                    if let httpResponse = response as? HTTPURLResponse {
                        print(httpResponse)
                        let t = (httpResponse.allHeaderFields["Token"] as? String)!
                        let Token = UserDefaults.standard
                        Token.set(t, forKey: "Token")
                        Token.synchronize()
                        DispatchQueue.main.async {
                            self.errorLabel.textColor = UIColor.systemGreen
                            self.errorLabel.text = "пароль успешно изменен"
                        }
                    }
                    
                } else {
                    print(httpResp.error_text)
                    DispatchQueue.main.async {
                        self.errorLabel.textColor = UIColor.systemRed
                        self.errorLabel.text = httpResp.error_text
                    }
                }
            } catch {
                print(error)
            }
        }.resume()
    }
    
    @IBAction func exitButton(_ sender: UIButton) {
        let Token = UserDefaults.standard
        Token.set("", forKey: "Token")
        Token.synchronize()
        tabBarController?.selectedIndex = 0
        
        let storyBoard = UIStoryboard(name: "Main", bundle: nil)
        let newVC = storyBoard.instantiateViewController(withIdentifier: "LogInViewController") as! LogInViewController
        self.present(newVC, animated: true, completion: nil)
    }

}
